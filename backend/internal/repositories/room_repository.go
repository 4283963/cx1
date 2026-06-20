package repositories

import (
	"smart-home/internal/database"
	"smart-home/internal/models"

	"github.com/google/uuid"
)

type RoomRepository struct{}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{}
}

func (r *RoomRepository) GetAll() ([]models.Room, error) {
	var rooms []models.Room
	err := database.DB.Find(&rooms).Error
	return rooms, err
}

func (r *RoomRepository) GetByID(id uuid.UUID) (*models.Room, error) {
	var room models.Room
	err := database.DB.First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) Create(room *models.Room) error {
	return database.DB.Create(room).Error
}

func (r *RoomRepository) Update(room *models.Room) error {
	return database.DB.Save(room).Error
}

func (r *RoomRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&models.Room{}, "id = ?", id).Error
}
