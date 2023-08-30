package segmentationService

import (
	"fmt"
	"regexp"
)

type SegmentPattern struct {
	Id      int    `json:"id"`
	Segment string `json:"segment" binding:"required"`
	Percent int    `json:"percent"`
}
type UpdateSegment struct {
	LastName string `json:"last_name" binding:"required"`
	NewName  string `json:"new_name" binding:"required"`
}

type Segmentation struct {
	Id               int        `json:"id"`
	Status           string     `json:"status"`
	AddingSegments   []string   `json:"adding_segments" binding:"-"`
	TTL              TimeToLive `json:"time_limit" binding:"-"`
	DeletingSegments []string   `json:"deleting_segments" binding:"-"`
	UserId           int        `json:"userId" binding:"required"`
}

type TimeToLive struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type ShowHistory struct {
	UserId    int    `json:"user_id" db:"user_id"`
	Segment   string `json:"segment" db:"segment"`
	Status    string `json:"status" db:"status"`
	Timestamp string `json:"timestamp" db:"date" binding:"-"`
}

func SlugValidation(segment string) error {
	regexPattern := "^[a-zA-Z0-9-]+$"
	matched, err := regexp.MatchString(regexPattern, segment)
	if err != nil {
		return fmt.Errorf("Incorrect pattern")
	}

	if !matched {
		return fmt.Errorf("Segment does not match slug pattern")
	}

	return nil
}
