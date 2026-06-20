package repositories

import (
	"smart-home/internal/database"
	"smart-home/internal/models"

	"github.com/google/uuid"
)

type LinkageRuleRepository struct{}

func NewLinkageRuleRepository() *LinkageRuleRepository {
	return &LinkageRuleRepository{}
}

func (r *LinkageRuleRepository) GetAll(roomID *uuid.UUID) ([]models.LinkageRule, error) {
	var rules []models.LinkageRule
	query := database.DB
	if roomID != nil {
		query = query.Where("room_id = ?", *roomID)
	}
	err := query.Order("created_at DESC").Find(&rules).Error
	return rules, err
}

func (r *LinkageRuleRepository) GetByID(id uuid.UUID) (*models.LinkageRule, error) {
	var rule models.LinkageRule
	err := database.DB.First(&rule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *LinkageRuleRepository) Create(rule *models.LinkageRule) error {
	return database.DB.Create(rule).Error
}

func (r *LinkageRuleRepository) Update(rule *models.LinkageRule) error {
	return database.DB.Save(rule).Error
}

func (r *LinkageRuleRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&models.LinkageRule{}, "id = ?", id).Error
}
