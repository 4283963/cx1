package repositories

import (
	"smart-home/internal/database"
	"smart-home/internal/models"
	"time"

	"github.com/google/uuid"
)

type EnvironmentRepository struct{}

func NewEnvironmentRepository() *EnvironmentRepository {
	return &EnvironmentRepository{}
}

func (r *EnvironmentRepository) Create(data *models.EnvironmentData) error {
	return database.DB.Create(data).Error
}

func (r *EnvironmentRepository) GetLatestByRoom(roomID uuid.UUID) (*models.EnvironmentData, error) {
	var data models.EnvironmentData
	err := database.DB.Where("room_id = ?", roomID).Order("timestamp DESC").First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *EnvironmentRepository) GetLatestAllRooms() ([]models.EnvironmentData, error) {
	subQuery := database.DB.Model(&models.EnvironmentData{}).
		Select("room_id, MAX(timestamp) as max_timestamp").
		Group("room_id")

	var data []models.EnvironmentData
	err := database.DB.Joins("JOIN (?) as sq ON environment_data.room_id = sq.room_id AND environment_data.timestamp = sq.max_timestamp", subQuery).
		Find(&data).Error
	return data, err
}

func (r *EnvironmentRepository) GetHistory(roomID uuid.UUID, startTime, endTime time.Time, limit int) ([]models.EnvironmentData, error) {
	var data []models.EnvironmentData
	query := database.DB.Where("room_id = ? AND timestamp BETWEEN ? AND ?", roomID, startTime, endTime)
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("timestamp ASC").Find(&data).Error
	return data, err
}

func (r *EnvironmentRepository) GetLatestByRoomID(roomID uuid.UUID) (*models.EnvironmentData, error) {
	return r.GetLatestByRoom(roomID)
}
