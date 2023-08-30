package service

import (
	segmentation_service "segmentationService"
	"segmentationService/src/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type User interface {
	CreateUser(user segmentation_service.User) (int, error)
	DeleteUser(user segmentation_service.User) error
	ShowUserSegments(userId int) (string, error)
	ShowHistory(userID int, startDate, endDate string) ([]segmentation_service.ShowHistory, error)
	CreateLinkToCSV(userID int, period, startDate, endDate string) (string, error)
}

type Segment interface {
	CreateSegment(segment segmentation_service.SegmentPattern) (int, error)
	DeleteSegment(pattern segmentation_service.SegmentPattern, input segmentation_service.Segmentation) error
	UpdateSegment(pattern segmentation_service.UpdateSegment) error
}

type Segmentation interface {
	SegmentMembership(input segmentation_service.Segmentation) error
	SegmentChecker(input segmentation_service.Segmentation) error
	AutoSegmentation(percent int, segment string, input segmentation_service.Segmentation) error
}

type Service struct {
	User
	Segment
	Segmentation
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		User:         NewUserSevice(repository.User),
		Segment:      NewSegmentSevice(repository.Segment, repository.Segmentation),
		Segmentation: NewSegmentationService(repository.Segmentation),
	}
}
