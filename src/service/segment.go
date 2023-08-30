package service

import (
	segmentation_service "segmentationService"
	"segmentationService/src/repository"
)

type SegmentSevice struct {
	repo             repository.Segment
	repoSegmentation repository.Segmentation
}

func NewSegmentSevice(repo repository.Segment, repoSegmentation repository.Segmentation) *SegmentSevice {
	return &SegmentSevice{repo: repo, repoSegmentation: repoSegmentation}
}

func (s *SegmentSevice) CreateSegment(pattern segmentation_service.SegmentPattern) (int, error) {
	return s.repo.CreateSegment(pattern)
}
func (s *SegmentSevice) DeleteSegment(pattern segmentation_service.SegmentPattern, input segmentation_service.Segmentation) error {
	users, err := s.repo.DeleteSegment(pattern)
	if err != nil {
		return err
	}
	for _, user := range users {
		input.UserId = user
		input.DeletingSegments = []string{pattern.Segment}
		input.AddingSegments = []string{}
		err = s.repoSegmentation.SegmentMembership(input)
		if err != nil {
			return err
		}
		err = s.repoSegmentation.SegmentationHistoryComment("deleted by admin", input.DeletingSegments, input.UserId)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *SegmentSevice) UpdateSegment(pattern segmentation_service.UpdateSegment) error {
	return s.repo.UpdateSegment(pattern)
}
