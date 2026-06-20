package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Floor     int       `json:"floor"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Device struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID     uuid.UUID `gorm:"type:uuid;not null" json:"room_id"`
	Name       string    `gorm:"not null" json:"name"`
	Type       string    `gorm:"not null" json:"type"`
	Status     bool      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Sensor struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID    uuid.UUID `gorm:"type:uuid;not null" json:"room_id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null" json:"type"`
	Unit      string    `json:"unit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LinkageRule struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID       uuid.UUID      `gorm:"type:uuid;not null" json:"room_id"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `json:"description"`
	TriggerType  string         `gorm:"not null" json:"trigger_type"`
	TriggerValue string         `gorm:"not null" json:"trigger_value"`
	ActionType   string         `gorm:"not null" json:"action_type"`
	ActionValue  string         `gorm:"not null" json:"action_value"`
	Enabled      bool           `json:"enabled"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type EnvironmentData struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID    uuid.UUID `gorm:"type:uuid;not null;index:idx_room_timestamp" json:"room_id"`
	Timestamp time.Time `gorm:"index:idx_room_timestamp" json:"timestamp"`
	Temp      float64   `json:"temp"`
	Humidity  float64   `json:"humidity"`
	Pm25      float64   `json:"pm25"`
	Formaldehyde float64 `json:"formaldehyde"`
}

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
