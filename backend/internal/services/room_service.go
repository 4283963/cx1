package services

import (
	"smart-home/internal/models"
	"smart-home/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type RoomService struct {
	repo *repositories.RoomRepository
}

func NewRoomService() *RoomService {
	return &RoomService{
		repo: repositories.NewRoomRepository(),
	}
}

func (s *RoomService) GetAll() ([]models.Room, error) {
	return s.repo.GetAll()
}

func (s *RoomService) GetByID(id uuid.UUID) (*models.Room, error) {
	return s.repo.GetByID(id)
}

func (s *RoomService) Create(name string, floor int) (*models.Room, error) {
	room := &models.Room{
		ID:        uuid.New(),
		Name:      name,
		Floor:     floor,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.Create(room)
	return room, err
}

func (s *RoomService) Update(id uuid.UUID, name string, floor int) (*models.Room, error) {
	room, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	room.Name = name
	room.Floor = floor
	room.UpdatedAt = time.Now()
	err = s.repo.Update(room)
	return room, err
}

func (s *RoomService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
