package handlers

import (
	"net/http"
	"smart-home/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LinkageRuleHandler struct {
	service *services.LinkageRuleService
}

func NewLinkageRuleHandler() *LinkageRuleHandler {
	return &LinkageRuleHandler{
		service: services.NewLinkageRuleService(),
	}
}

type CreateLinkageRuleRequest struct {
	RoomID       uuid.UUID `json:"room_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description"`
	TriggerType  string    `json:"trigger_type" binding:"required"`
	TriggerValue string    `json:"trigger_value" binding:"required"`
	ActionType   string    `json:"action_type" binding:"required"`
	ActionValue  string    `json:"action_value" binding:"required"`
	Enabled      bool      `json:"enabled"`
}

type UpdateLinkageRuleRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	TriggerType  string `json:"trigger_type" binding:"required"`
	TriggerValue string `json:"trigger_value" binding:"required"`
	ActionType   string `json:"action_type" binding:"required"`
	ActionValue  string `json:"action_value" binding:"required"`
	Enabled      bool   `json:"enabled"`
}

func (h *LinkageRuleHandler) GetAll(c *gin.Context) {
	var roomID *uuid.UUID
	roomIDStr := c.Query("room_id")
	if roomIDStr != "" {
		id, err := uuid.Parse(roomIDStr)
		if err == nil {
			roomID = &id
		}
	}
	rules, err := h.service.GetAll(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

func (h *LinkageRuleHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}
	rule, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rule})
}

func (h *LinkageRuleHandler) Create(c *gin.Context) {
	var req CreateLinkageRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule, err := h.service.Create(req.RoomID, req.Name, req.Description, req.TriggerType, req.TriggerValue, req.ActionType, req.ActionValue, req.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": rule})
}

func (h *LinkageRuleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}
	var req UpdateLinkageRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule, err := h.service.Update(id, req.Name, req.Description, req.TriggerType, req.TriggerValue, req.ActionType, req.ActionValue, req.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rule})
}

func (h *LinkageRuleHandler) ToggleEnabled(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}
	rule, err := h.service.ToggleEnabled(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rule})
}

func (h *LinkageRuleHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rule deleted successfully"})
}
