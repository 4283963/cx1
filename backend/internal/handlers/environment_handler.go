package handlers

import (
	"net/http"
	"smart-home/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EnvironmentHandler struct {
	service *services.EnvironmentService
}

func NewEnvironmentHandler() *EnvironmentHandler {
	return &EnvironmentHandler{
		service: services.NewEnvironmentService(),
	}
}

func (h *EnvironmentHandler) GetLatestAll(c *gin.Context) {
	data, err := h.service.GetLatestAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *EnvironmentHandler) GetLatestByRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	data, err := h.service.GetLatestByRoom(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment data not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *EnvironmentHandler) GetHistory(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	limitStr := c.Query("limit")

	var startTime, endTime time.Time
	var limit int

	if startTimeStr != "" {
		startTime, _ = time.Parse(time.RFC3339, startTimeStr)
	} else {
		startTime = time.Now().Add(-24 * time.Hour)
	}

	if endTimeStr != "" {
		endTime, _ = time.Parse(time.RFC3339, endTimeStr)
	} else {
		endTime = time.Now()
	}

	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	data, err := h.service.GetHistory(roomID, startTime, endTime, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
