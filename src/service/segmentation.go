package service

import (
	"errors"
	"math"
	"math/rand"
	segmentation_service "segmentationService"
	"segmentationService/src/repository"
	"time"
)

type SegmentationService struct {
	repo repository.Segmentation
}

func NewSegmentationService(repo repository.Segmentation) *SegmentationService {
	return &SegmentationService{repo: repo}
}

func (s *SegmentationService) SegmentMembership(input segmentation_service.Segmentation) error {
	if len(input.AddingSegments) == 0 && len(input.DeletingSegments) == 0 {
		return errors.New("Add to JSON AddingSegments or DeletingSegments")
	}
	if input.TTL.Value != 0 && input.TTL.Unit != "" {
		TTL, _ := s.TimeParser(input)
		err := s.repo.SegmentMembership(input)
		if err != nil {
			return err
		}
		err = s.repo.SegmentationHistoryComment("added", input.AddingSegments, input.UserId)
		if err != nil {
			return err
		}
		go func() error {
			for {
				currentTime := time.Now().UTC()
				if currentTime.After(TTL) {
					input.DeletingSegments = input.AddingSegments
					input.AddingSegments = []string{}
					err = s.repo.SegmentMembership(input)
					if err != nil {
						return err
					}
					err = s.repo.SegmentationHistoryComment("time limited", input.DeletingSegments, input.UserId)
					if err != nil {
						return err
					}
					break
				}
				time.Sleep(time.Second)
			}
			return err
		}()
	} else {
		err := s.repo.SegmentationHistoryComment("added", input.AddingSegments, input.UserId)
		if err != nil {
			return err
		}
		err = s.repo.SegmentationHistoryComment("deleted", input.DeletingSegments, input.UserId)
		if err != nil {
			return err
		}
		err = s.repo.SegmentMembership(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SegmentationService) SegmentChecker(input segmentation_service.Segmentation) error {
	return s.repo.SegmentChecker(input)
}

func (s *SegmentationService) TimeParser(input segmentation_service.Segmentation) (time.Time, error) {
	valueDuration := time.Duration(input.TTL.Value)
	value := input.TTL.Value

	var TTL time.Time
	switch input.TTL.Unit {
	case "h":
		TTL = time.Now().Add(time.Hour * valueDuration)
	case "min":
		TTL = time.Now().Add(time.Minute * valueDuration)
	case "second":
		TTL = time.Now().Add(time.Second * valueDuration)
	case "days":
		TTL = time.Now().AddDate(0, 0, value)
	case "month":
		TTL = time.Now().AddDate(0, value, 0)
	case "years":
		TTL = time.Now().AddDate(value, 0, 0)
	}
	return TTL, nil
}

func (s *SegmentationService) UsersForAutoSegmentation(percent int) ([]int, error) {
	usersCountInTable, _ := s.repo.CountOfUsers()
	generalCountOfUsersWithSegment := int(math.Floor(float64(usersCountInTable * percent / 100)))

	usersWithSegment := make([]int, generalCountOfUsersWithSegment)
	for i := 0; i < generalCountOfUsersWithSegment; i++ {
		usersWithSegment[i] = i + 1
	}

	shuffle(usersWithSegment)

	return usersWithSegment, nil
}

func (s *SegmentationService) AutoSegmentation(percent int, segment string, input segmentation_service.Segmentation) error {
	input.AddingSegments = []string{segment}
	usersWithSegment, _ := s.UsersForAutoSegmentation(percent)

	for _, user := range usersWithSegment {
		input.UserId = user
		err := s.repo.SegmentationHistoryComment("auto added", input.AddingSegments, input.UserId)
		if err != nil {
			return err
		}
		err = s.repo.SegmentMembership(input)
		if err != nil {
			return err
		}
	}

	return nil
}

func shuffle(arr []int) {
	rand.Seed(time.Now().UnixNano())
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}
