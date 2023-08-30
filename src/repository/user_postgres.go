package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	segmentation_service "segmentationService"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user segmentation_service.User) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var id int
	createUserQuery := fmt.Sprintf("INSERT INTO %s (username) VALUES ($1) RETURNING id", usersTable)
	if err := tx.QueryRow(createUserQuery, user.Username).Scan(&id); err != nil {
		return 0, err
	}

	addUserToSegmentsInfoQuery := fmt.Sprintf("INSERT INTO %s (user_id, segments_list) VALUES ($1, $2)", segmentationInfo)
	_, err = tx.Exec(addUserToSegmentsInfoQuery, id, "EMPTY")
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (r *UserRepository) DeleteUser(user segmentation_service.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	deleteUser := fmt.Sprintf("DELETE FROM %s WHERE username = $1 RETURNING id", usersTable)
	var deletedUserId int
	err = tx.QueryRow(deleteUser, user.Username).Scan(&deletedUserId)
	if err != nil {
		return err
	}
	deleteUserFromSegmentationInfo := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", segmentationInfo)
	_, err = tx.Exec(deleteUserFromSegmentationInfo, deletedUserId)
	if err != nil {
		return err
	}

	deleteUserFromSegmentationHistory := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", segmentationHistory)
	_, err = tx.Exec(deleteUserFromSegmentationHistory, deletedUserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) ShowHistory(userID int, startDate, endDate string) ([]segmentation_service.ShowHistory, error) {
	var history []segmentation_service.ShowHistory
	query := fmt.Sprintf("SELECT user_id, segment, status, date FROM %s WHERE user_id=$1 AND date >= '%s' AND date < '%s'", segmentationHistory, startDate, endDate)
	err := r.db.Select(&history, query, userID)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (r *UserRepository) ShowUserSegments(userId int) (string, error) {
	userSegments := fmt.Sprintf("SELECT segments_list FROM %s WHERE user_id = $1", segmentationInfo)
	var segmentsList string
	err := r.db.QueryRow(userSegments, userId).Scan(&segmentsList)
	if err != nil {
		return "", err
	}
	return segmentsList, err
}
