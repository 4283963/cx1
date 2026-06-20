package handlers

import (
	"net/http"
	"smart-home/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeviceHandler struct {
	service *services.DeviceService
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{
		service: services.NewDeviceService(),
	}
}

func (h *DeviceHandler) GetAll(c *gin.Context) {
	var roomID *uuid.UUID
	roomIDStr := c.Query("room_id")
	if roomIDStr != "" {
		id, err := uuid.Parse(roomIDStr)
		if err == nil {
			roomID = &id
		}
	}
	devices, err := h.service.GetAll(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": devices})
}

func (h *DeviceHandler) ToggleStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}
	device, err := h.service.ToggleStatus(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": device})
}
