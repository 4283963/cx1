package main

import (
	"log"
	"smart-home/internal/config"
	"smart-home/internal/database"
	"smart-home/internal/gateway"
	"smart-home/internal/handlers"
	"smart-home/internal/middleware"
	"smart-home/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.InitDB()
	websocket.StartHub()

	go func() {
		simulator := gateway.NewGatewaySimulator()
		simulator.Start()
	}()

	r := gin.Default()
	r.Use(middleware.CORS())

	roomHandler := handlers.NewRoomHandler()
	deviceHandler := handlers.NewDeviceHandler()
	linkageRuleHandler := handlers.NewLinkageRuleHandler()
	environmentHandler := handlers.NewEnvironmentHandler()

	api := r.Group("/api/v1")
	{
		rooms := api.Group("/rooms")
		{
			rooms.GET("", roomHandler.GetAll)
			rooms.GET("/:id", roomHandler.GetByID)
			rooms.POST("", roomHandler.Create)
			rooms.PUT("/:id", roomHandler.Update)
			rooms.DELETE("/:id", roomHandler.Delete)
		}

		devices := api.Group("/devices")
		{
			devices.GET("", deviceHandler.GetAll)
			devices.PATCH("/:id/toggle", deviceHandler.ToggleStatus)
		}

		rules := api.Group("/linkage-rules")
		{
			rules.GET("", linkageRuleHandler.GetAll)
			rules.GET("/:id", linkageRuleHandler.GetByID)
			rules.POST("", linkageRuleHandler.Create)
			rules.PUT("/:id", linkageRuleHandler.Update)
			rules.PATCH("/:id/toggle", linkageRuleHandler.ToggleEnabled)
			rules.DELETE("/:id", linkageRuleHandler.Delete)
		}

		environment := api.Group("/environment")
		{
			environment.GET("/latest", environmentHandler.GetLatestAll)
			environment.GET("/latest/:room_id", environmentHandler.GetLatestByRoom)
			environment.GET("/history/:room_id", environmentHandler.GetHistory)
		}
	}

	r.GET("/ws", websocket.HandleWebSocket)

	log.Printf("Server starting on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
