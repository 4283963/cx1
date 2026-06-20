package handlers

import (
	"net/http"
	"smart-home/internal/services"
	"smart-home/internal/websocket"
	"smart-home/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	service *services.SystemService
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		service: services.NewSystemService(),
	}
}

type SetForceModeRequest struct {
	Enabled utils.FlexBool `json:"enabled" binding:"required"`
}

func (h *SystemHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"force_mode": h.service.GetForceMode(),
		},
	})
}

func (h *SystemHandler) SetForceMode(c *gin.Context) {
	var req SetForceModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.service.SetForceMode(req.Enabled.Bool())

	websocket.BroadcastMessage("force_mode_update", gin.H{
		"enabled": req.Enabled.Bool(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Force mode updated",
		"data": gin.H{
			"force_mode": req.Enabled.Bool(),
		},
	})
}
