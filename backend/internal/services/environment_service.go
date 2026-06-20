package services

import (
	"smart-home/internal/models"
	"smart-home/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type EnvironmentService struct {
	repo *repositories.EnvironmentRepository
}

func NewEnvironmentService() *EnvironmentService {
	return &EnvironmentService{
		repo: repositories.NewEnvironmentRepository(),
	}
}

func (s *EnvironmentService) Create(roomID uuid.UUID, temp, humidity, pm25, formaldehyde float64) (*models.EnvironmentData, error) {
	data := &models.EnvironmentData{
		ID:           uuid.New(),
		RoomID:       roomID,
		Timestamp:    time.Now(),
		Temp:         temp,
		Humidity:     humidity,
		Pm25:         pm25,
		Formaldehyde: formaldehyde,
	}
	err := s.repo.Create(data)
	return data, err
}

func (s *EnvironmentService) GetLatestByRoom(roomID uuid.UUID) (*models.EnvironmentData, error) {
	return s.repo.GetLatestByRoom(roomID)
}

func (s *EnvironmentService) GetLatestAllRooms() ([]models.EnvironmentData, error) {
	return s.repo.GetLatestAllRooms()
}

func (s *EnvironmentService) GetHistory(roomID uuid.UUID, startTime, endTime time.Time, limit int) ([]models.EnvironmentData, error) {
	return s.repo.GetHistory(roomID, startTime, endTime, limit)
}
