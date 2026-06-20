package services

import (
	"smart-home/internal/models"
	"smart-home/internal/repositories"

	"github.com/google/uuid"
)

type DeviceService struct {
	repo *repositories.DeviceRepository
}

func NewDeviceService() *DeviceService {
	return &DeviceService{
		repo: repositories.NewDeviceRepository(),
	}
}

func (s *DeviceService) GetAll(roomID *uuid.UUID) ([]models.Device, error) {
	return s.repo.GetAll(roomID)
}

func (s *DeviceService) GetByID(id uuid.UUID) (*models.Device, error) {
	return s.repo.GetByID(id)
}

func (s *DeviceService) ToggleStatus(id uuid.UUID) (*models.Device, error) {
	device, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	device.Status = !device.Status
	err = s.repo.Update(device)
	return device, err
}
