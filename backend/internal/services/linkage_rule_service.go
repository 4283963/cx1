package services

import (
	"smart-home/internal/models"
	"smart-home/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type LinkageRuleService struct {
	repo *repositories.LinkageRuleRepository
}

func NewLinkageRuleService() *LinkageRuleService {
	return &LinkageRuleService{
		repo: repositories.NewLinkageRuleRepository(),
	}
}

func (s *LinkageRuleService) GetAll(roomID *uuid.UUID) ([]models.LinkageRule, error) {
	return s.repo.GetAll(roomID)
}

func (s *LinkageRuleService) GetByID(id uuid.UUID) (*models.LinkageRule, error) {
	return s.repo.GetByID(id)
}

func (s *LinkageRuleService) Create(roomID uuid.UUID, name, description, triggerType, triggerValue, actionType, actionValue string, enabled bool) (*models.LinkageRule, error) {
	rule := &models.LinkageRule{
		ID:           uuid.New(),
		RoomID:       roomID,
		Name:         name,
		Description:  description,
		TriggerType:  triggerType,
		TriggerValue: triggerValue,
		ActionType:   actionType,
		ActionValue:  actionValue,
		Enabled:      enabled,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := s.repo.Create(rule)
	return rule, err
}

func (s *LinkageRuleService) Update(id uuid.UUID, name, description, triggerType, triggerValue, actionType, actionValue string, enabled bool) (*models.LinkageRule, error) {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	rule.Name = name
	rule.Description = description
	rule.TriggerType = triggerType
	rule.TriggerValue = triggerValue
	rule.ActionType = actionType
	rule.ActionValue = actionValue
	rule.Enabled = enabled
	rule.UpdatedAt = time.Now()
	err = s.repo.Update(rule)
	return rule, err
}

func (s *LinkageRuleService) ToggleEnabled(id uuid.UUID) (*models.LinkageRule, error) {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	rule.Enabled = !rule.Enabled
	rule.UpdatedAt = time.Now()
	err = s.repo.Update(rule)
	return rule, err
}

func (s *LinkageRuleService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
