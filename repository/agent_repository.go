package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type AgentRepository struct{}

func NewAgentRepository() *AgentRepository { return &AgentRepository{} }

func (r *AgentRepository) GetAll() ([]models.Agent, error) {
	var rows []models.Agent
	err := config.DB.Find(&rows).Error
	return rows, err
}

func (r *AgentRepository) GetByUserID(userID uint) (*models.Agent, error) {
	var row models.Agent
	if err := config.DB.First(&row, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *AgentRepository) Create(a *models.Agent) error {
	return config.DB.Create(a).Error
}

func (r *AgentRepository) Delete(userID uint) error {
	return config.DB.Delete(&models.Agent{}, "user_id = ?", userID).Error
}
