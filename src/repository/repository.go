package repository

import (
	"github.com/jmoiron/sqlx"
	segmentation_service "segmentationService"
)

type User interface {
	CreateUser(user segmentation_service.User) (int, error)
	DeleteUser(user segmentation_service.User) error
	ShowUserSegments(userId int) (string, error)
	ShowHistory(userID int, startDate, endDate string) ([]segmentation_service.ShowHistory, error)
}

type Segment interface {
	CreateSegment(pattern segmentation_service.SegmentPattern) (int, error)
	DeleteSegment(pattern segmentation_service.SegmentPattern) ([]int, error)
	UpdateSegment(pattern segmentation_service.UpdateSegment) error
}

type Segmentation interface {
	SegmentMembership(input segmentation_service.Segmentation) error
	SegmentChecker(input segmentation_service.Segmentation) error
	CountOfUsers() (int, []int, error)
	SegmentationHistoryComment(status string, segments []string, userId int) error
}

type Repository struct {
	User
	Segment
	Segmentation
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:         NewUserRepository(db),
		Segment:      NewSegmentRepository(db),
		Segmentation: NewSegmentationRepository(db),
	}
}
