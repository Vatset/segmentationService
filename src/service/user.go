package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	segmentation_service "segmentationService"
	"segmentationService/src/repository"
	"strconv"
)

type UserSevice struct {
	repo repository.User
}

func NewUserSevice(repo repository.User) *UserSevice {
	return &UserSevice{repo: repo}
}

func (s *UserSevice) CreateUser(user segmentation_service.User) (int, error) {
	return s.repo.CreateUser(user)
}

func (s *UserSevice) DeleteUser(user segmentation_service.User) error {
	return s.repo.DeleteUser(user)
}

func (s *UserSevice) CreateLinkToCSV(userID int, period, startDate, endDate string) (string, error) {
	history, err := s.repo.ShowHistory(userID, startDate, endDate)
	if err != nil {
		return "", errors.New("error with getting history")

	}
	filename := "user" + strconv.Itoa(userID) + "_period" + period + "_history.csv"

	folderPath := "history"
	rootPath, err := os.Getwd()
	if err != nil {
		return "", errors.New("Error getting directory")
	}

	fullPath := filepath.Join(rootPath, folderPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.Mkdir(fullPath, 0755); err != nil {
			return "", errors.New("Error creating folder")
		}
	}

	savePath := "history/"
	filePath := filepath.Join(savePath, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("error with creating file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"User ID", "Segment", "Status", "Timestamp"})

	for _, historyEntry := range history {
		row := []string{
			fmt.Sprintf("%d", historyEntry.UserId),
			historyEntry.Segment,
			historyEntry.Status,
			historyEntry.Timestamp,
		}
		writer.Write(row)
	}

	link := "http://localhost:8080/api/user/history/" + filename
	return link, nil
}

func (s *UserSevice) ShowUserSegments(userId int) (string, error) {
	return s.repo.ShowUserSegments(userId)
}
func (s *UserSevice) ShowHistory(userID int, startDate, endDate string) ([]segmentation_service.ShowHistory, error) {
	return s.repo.ShowHistory(userID, startDate, endDate)
}
