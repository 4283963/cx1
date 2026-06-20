package repositories

import (
	"smart-home/internal/database"
	"smart-home/internal/models"

	"github.com/google/uuid"
)

type DeviceRepository struct{}

func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{}
}

func (r *DeviceRepository) GetAll(roomID *uuid.UUID) ([]models.Device, error) {
	var devices []models.Device
	query := database.DB
	if roomID != nil {
		query = query.Where("room_id = ?", *roomID)
	}
	err := query.Find(&devices).Error
	return devices, err
}

func (r *DeviceRepository) GetByID(id uuid.UUID) (*models.Device, error) {
	var device models.Device
	err := database.DB.First(&device, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) Update(device *models.Device) error {
	return database.DB.Save(device).Error
}
