package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	segmentation_service "segmentationService"
	"strings"
	"time"
)

type SegmentationRepository struct {
	db *sqlx.DB
}

func NewSegmentationRepository(db *sqlx.DB) *SegmentationRepository {
	return &SegmentationRepository{db: db}
}

func SetCurrentTime() time.Time {
	date := time.Now()
	return date
}

func (r *SegmentationRepository) SegmentChecker(input segmentation_service.Segmentation) error {
	for _, segment := range input.AddingSegments {
		checkSegmentExist := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE segment = $1", segmentsTable)
		var count int
		err := r.db.QueryRow(checkSegmentExist, segment).Scan(&count)
		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Segment that you trying to add %s does not exist", segment)
		}
	}

	for _, segment := range input.DeletingSegments {
		checkSegmentExist := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE segment = $1", segmentsTable)
		var count int
		err := r.db.QueryRow(checkSegmentExist, segment).Scan(&count)
		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Segment that you trying to delete %s does not exist", segment)
		}

	}
	for _, segment := range input.AddingSegments {
		segmentationsList, _ := r.UserSegments(input)
		if strings.Contains(segmentationsList, segment) {
			return fmt.Errorf("Segment that you are trying to add '%s' already exists", segment)
		}
	}
	if input.UserId != 0 {
		for _, segment := range input.DeletingSegments {
			segmentationsList, _ := r.UserSegments(input)
			if !strings.Contains(segmentationsList, segment) {
				return fmt.Errorf("The user does not belong to the segment '%s' you are trying to delete", segment)
			}
		}

	}

	return nil
}

func (r *SegmentationRepository) UserSegments(input segmentation_service.Segmentation) (string, error) {
	userSegments := fmt.Sprintf("SELECT segments_list FROM %s WHERE user_id = $1", segmentationInfo)
	var segmentationsList string
	err := r.db.QueryRow(userSegments, input.UserId).Scan(&segmentationsList)
	if err != nil {
		return "", errors.New("user does not exist")
	}
	return segmentationsList, err
}

func (r *SegmentationRepository) AddSegmentation(segmentationsList string, input segmentation_service.Segmentation) error {
	if len(input.AddingSegments) != 0 {
		addingSegmentsStr := strings.Join(input.AddingSegments, ",")

		if segmentationsList == "EMPTY" {
			addUserSegments := fmt.Sprintf("UPDATE %s SET segments_list = $2 WHERE user_id=$1", segmentationInfo)
			_, err := r.db.Exec(addUserSegments, input.UserId, addingSegmentsStr)
			return err
		} else {
			newSegmentationsList := segmentationsList + "," + addingSegmentsStr
			addUserSegments := fmt.Sprintf("UPDATE %s SET segments_list = $2 WHERE user_id=$1", segmentationInfo)
			_, err := r.db.Exec(addUserSegments, input.UserId, newSegmentationsList)
			return err
		}

	}
	return nil
}
func (r *SegmentationRepository) DeleteSegmentation(segmentationsList string, input segmentation_service.Segmentation) error {
	if len(input.DeletingSegments) != 0 {
		segmentationsListArr := strings.Split(segmentationsList, ",")
		var newSegmentationsListArr []string

		for _, segment := range segmentationsListArr {
			found := false
			for _, deletingSegment := range input.DeletingSegments {
				if segment == deletingSegment {
					found = true
					break
				}
			}
			if !found {
				newSegmentationsListArr = append(newSegmentationsListArr, segment)
			}
		}
		var newSegmentationsList string
		if len(newSegmentationsListArr) == 0 {
			newSegmentationsList = "EMPTY"
		} else {
			newSegmentationsList = strings.Join(newSegmentationsListArr, ",")
		}

		deleteUserSegments := fmt.Sprintf("UPDATE %s SET segments_list = $2 WHERE user_id=$1", segmentationInfo)
		_, err := r.db.Exec(deleteUserSegments, input.UserId, newSegmentationsList)
		if err != nil {
			return err
		}

		return err
	}
	return nil
}

func (r *SegmentationRepository) SegmentMembership(input segmentation_service.Segmentation) error {
	segmentationsList, err := r.UserSegments(input)
	if err != nil {
		return err
	}
	err = r.AddSegmentation(segmentationsList, input)
	if err != nil {
		errors.New("error with adding segmentation")
	}

	err = r.DeleteSegmentation(segmentationsList, input)
	if err != nil {
		errors.New("error with deleting segmentation")
	}
	return err
}
func (r *SegmentationRepository) SegmentationHistoryComment(status string, segments []string, userId int) error {
	userSegmentHistory := fmt.Sprintf("INSERT INTO %s (user_id, segment, status, date) VALUES ($1, $2, $3, $4)", segmentationHistory)

	for _, segment := range segments {
		_, err := r.db.Exec(userSegmentHistory, userId, segment, status, SetCurrentTime())
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *SegmentationRepository) CountOfUsers() (int, []int, error) {
	query := fmt.Sprintf("SELECT id FROM %s", usersTable)
	rows, err := r.db.Query(query)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var users []int
	var user int
	for rows.Next() {
		if err := rows.Scan(&user); err != nil {
			return 0, nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return 0, nil, err
	}

	return len(users), users, nil
}
