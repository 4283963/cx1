package database

import (
	"fmt"
	"log"
	"smart-home/internal/config"
	"smart-home/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.Room{},
		&models.Device{},
		&models.Sensor{},
		&models.LinkageRule{},
		&models.EnvironmentData{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedData()
	log.Println("Database initialized successfully")
}

func seedData() {
	var roomCount int64
	DB.Model(&models.Room{}).Count(&roomCount)
	if roomCount > 0 {
		return
	}

	rooms := []models.Room{
		{ID: uuid.New(), Name: "客厅", Floor: 1},
		{ID: uuid.New(), Name: "主卧", Floor: 2},
		{ID: uuid.New(), Name: "次卧", Floor: 2},
		{ID: uuid.New(), Name: "厨房", Floor: 1},
		{ID: uuid.New(), Name: "卫生间", Floor: 1},
	}

	for i := range rooms {
		DB.Create(&rooms[i])
	}

	for _, room := range rooms {
		devices := []models.Device{
			{ID: uuid.New(), RoomID: room.ID, Name: "智能灯", Type: "light", Status: false},
			{ID: uuid.New(), RoomID: room.ID, Name: "空调", Type: "ac", Status: false},
			{ID: uuid.New(), RoomID: room.ID, Name: "空气净化器", Type: "purifier", Status: false},
		}
		for i := range devices {
			DB.Create(&devices[i])
		}

		sensors := []models.Sensor{
			{ID: uuid.New(), RoomID: room.ID, Name: "温度传感器", Type: "temp", Unit: "°C"},
			{ID: uuid.New(), RoomID: room.ID, Name: "湿度传感器", Type: "humidity", Unit: "%"},
			{ID: uuid.New(), RoomID: room.ID, Name: "PM2.5传感器", Type: "pm25", Unit: "μg/m³"},
			{ID: uuid.New(), RoomID: room.ID, Name: "甲醛传感器", Type: "formaldehyde", Unit: "mg/m³"},
		}
		for i := range sensors {
			DB.Create(&sensors[i])
		}
	}

	log.Println("Seed data created successfully")
}
