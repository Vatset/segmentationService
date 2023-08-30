package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	segmentation_service "segmentationService"
	"strings"
)

type SegmentRepository struct {
	db *sqlx.DB
}

func NewSegmentRepository(db *sqlx.DB) *SegmentRepository {
	return &SegmentRepository{db: db}
}

func (r *SegmentRepository) CreateSegment(pattern segmentation_service.SegmentPattern) (int, error) {
	var id int
	createSegment := fmt.Sprintf("INSERT INTO %s (segment) VALUES ($1) RETURNING id", segmentsTable)
	row := r.db.QueryRow(createSegment, pattern.Segment)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SegmentRepository) DeleteSegment(pattern segmentation_service.SegmentPattern) ([]int, error) {
	deleteSegment := fmt.Sprintf("DELETE FROM %s WHERE segment = $1", segmentsTable)
	_, err := r.db.Exec(deleteSegment, pattern.Segment)
	if err != nil {
		return nil, err
	}

	segmentationsListTable := fmt.Sprintf("SELECT segments_list,user_id FROM %s", segmentationInfo)
	rows, err := r.db.Query(segmentationsListTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usersWithSegment := []int{}
	for rows.Next() {
		var segmentationsList string
		var user_id int

		if err := rows.Scan(&segmentationsList, &user_id); err != nil {
			return nil, err
		}
		segmentationsListArr := strings.Split(segmentationsList, ",")
		for _, segment := range segmentationsListArr {
			if segment == pattern.Segment {
				usersWithSegment = append(usersWithSegment, user_id)
			}

		}
	}
	return usersWithSegment, err
}

func (r *SegmentRepository) UpdateSegment(pattern segmentation_service.UpdateSegment) error {
	err := r.UpdateSegmentChecker(pattern)
	if err != nil {
		return err
	}

	updateSegment := fmt.Sprintf("UPDATE %s SET segment=$2 WHERE segment=$1 ", segmentsTable)
	_, err = r.db.Exec(updateSegment, pattern.LastName, pattern.NewName)
	if err != nil {
		errors.New("error with updating segment in segments table")
	}

	err = r.SegmentsInTableChanger(segmentationsListField, segmentationInfo, pattern.LastName, pattern.NewName)
	if err != nil {
		errors.New("error with updating segment in segmentation_info table")
	}
	err = r.SegmentsInTableChanger(segmentationsNameField, segmentationHistory, pattern.LastName, pattern.NewName)
	if err != nil {
		errors.New("error with updating segment in segmentations_history table")
	}

	return nil
}

func (r *SegmentRepository) UpdateSegmentChecker(input segmentation_service.UpdateSegment) error {
	checkSegmentExist := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE segment = $1", segmentsTable)
	var count int
	err := r.db.QueryRow(checkSegmentExist, input.LastName).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("Segment that you trying to update does not exist")
	}

	return nil
}
func (r *SegmentRepository) SegmentsInTableChanger(row, table, lastNameOfSegment, newNameofSegment string) error {
	segments := fmt.Sprintf("SELECT id, %s FROM %s", row, table)
	rows, err := r.db.Query(segments)
	if err != nil {
		return err
	}
	defer rows.Close()

	var updates []struct {
		ID              int
		NewSegmentsList string
	}

	for rows.Next() {
		var id int
		var segmentsList string
		if err := rows.Scan(&id, &segmentsList); err != nil {
			return err
		}
		segmentationsListArr := strings.Split(segmentsList, ",")
		var newsegmentsListArr []string
		for _, segment := range segmentationsListArr {
			if segment != lastNameOfSegment {
				newsegmentsListArr = append(newsegmentsListArr, segment)
			} else {
				newsegmentsListArr = append(newsegmentsListArr, newNameofSegment)
			}
		}
		newSegmentsList := strings.Join(newsegmentsListArr, ",")
		updates = append(updates, struct {
			ID              int
			NewSegmentsList string
		}{
			ID:              id,
			NewSegmentsList: newSegmentsList,
		})
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, update := range updates {
		updateQuery := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = $2", table, row)
		_, err = tx.Exec(updateQuery, update.NewSegmentsList, update.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
