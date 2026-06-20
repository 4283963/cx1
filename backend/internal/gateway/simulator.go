package gateway

import (
	"log"
	"math/rand"
	"smart-home/internal/config"
	"smart-home/internal/database"
	"smart-home/internal/models"
	"smart-home/internal/websocket"
	"time"

	"github.com/google/uuid"
)

type GatewaySimulator struct {
	roomData map[uuid.UUID]*RoomSensorData
}

type RoomSensorData struct {
	RoomID       uuid.UUID
	RoomName     string
	Temp         float64
	Humidity     float64
	Pm25         float64
	Formaldehyde float64
}

func NewGatewaySimulator() *GatewaySimulator {
	sim := &GatewaySimulator{
		roomData: make(map[uuid.UUID]*RoomSensorData),
	}
	sim.initRooms()
	return sim
}

func (s *GatewaySimulator) initRooms() {
	var rooms []models.Room
	result := database.DB.Find(&rooms)
	if result.Error != nil {
		log.Printf("Error loading rooms: %v", result.Error)
		return
	}

	baseTemp := 22.0
	baseHumidity := 50.0

	for _, room := range rooms {
		s.roomData[room.ID] = &RoomSensorData{
			RoomID:       room.ID,
			RoomName:     room.Name,
			Temp:         baseTemp + rand.Float64()*4 - 2,
			Humidity:     baseHumidity + rand.Float64()*10 - 5,
			Pm25:         rand.Float64() * 50,
			Formaldehyde: rand.Float64() * 0.1,
		}
	}
}

func (s *GatewaySimulator) Start() {
	interval := time.Duration(config.AppConfig.GatewaySimulationInterval) * time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Gateway simulator started, reporting interval: %v", interval)

	for range ticker.C {
		s.simulateDataReport()
	}
}

func (s *GatewaySimulator) simulateDataReport() {
	var allData []map[string]interface{}

	for roomID, data := range s.roomData {
		s.updateSensorData(data)

		envData := &models.EnvironmentData{
			ID:           uuid.New(),
			RoomID:       roomID,
			Timestamp:    time.Now(),
			Temp:         data.Temp,
			Humidity:     data.Humidity,
			Pm25:         data.Pm25,
			Formaldehyde: data.Formaldehyde,
		}

		if err := database.DB.Create(envData).Error; err != nil {
			log.Printf("Error saving environment data: %v", err)
			continue
		}

		allData = append(allData, map[string]interface{}{
			"room_id":       roomID.String(),
			"room_name":     data.RoomName,
			"timestamp":     envData.Timestamp,
			"temp":          data.Temp,
			"humidity":      data.Humidity,
			"pm25":          data.Pm25,
			"formaldehyde":  data.Formaldehyde,
			"temp_status":   getStatus(data.Temp, 18, 26, 30),
			"humidity_status": getStatus(data.Humidity, 40, 60, 70),
			"pm25_status":   getStatus(data.Pm25, 0, 35, 75),
			"formaldehyde_status": getStatus(data.Formaldehyde, 0, 0.08, 0.12),
		})
	}

	websocket.BroadcastMessage("environment_data", allData)
}

func (s *GatewaySimulator) updateSensorData(data *RoomSensorData) {
	data.Temp += (rand.Float64() - 0.5) * 0.3
	data.Temp = clamp(data.Temp, 10, 35)

	data.Humidity += (rand.Float64() - 0.5) * 1
	data.Humidity = clamp(data.Humidity, 20, 90)

	data.Pm25 += (rand.Float64() - 0.5) * 8
	data.Pm25 = clamp(data.Pm25, 0, 300)

	data.Formaldehyde += (rand.Float64() - 0.5) * 0.005
	data.Formaldehyde = clamp(data.Formaldehyde, 0, 0.3)
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func getStatus(value float64, goodMax, warningMax, dangerMax float64) string {
	if value <= goodMax {
		return "good"
	} else if value <= warningMax {
		return "warning"
	}
	return "danger"
}
