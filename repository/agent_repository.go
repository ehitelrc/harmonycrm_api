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

func (r *AgentRepository) GetAllByCompanyIDWithUserInfo(companyID uint) ([]models.AgentDepartmentInformation, error) {
	var rows []models.AgentDepartmentInformation
	err := config.DB.Where("company_id = ?", companyID).Find(&rows).Error
	return rows, err
}

func (r *AgentRepository) GetAllByCompanyIDAndDepartmentIDWithUserInfo(companyID uint, departmentID uint) ([]models.AgentDepartmentInformation, error) {
	var rows []models.AgentDepartmentInformation
	err := config.DB.Where("company_id = ? AND department_id = ?", companyID, departmentID).Find(&rows).Error
	return rows, err
}

func (r *AgentRepository) GetAllWithUserInfo() ([]models.AgentUser, error) {
	var rows []models.AgentUser
	err := config.DB.Find(&rows).Error
	return rows, err
}

func (r *AgentRepository) GetByUserIDWithUserInfo(userID uint) (*models.AgentUser, error) {
	var row models.AgentUser
	if err := config.DB.First(&row, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *AgentRepository) GetAllNonAgents() ([]models.NonAgentUser, error) {
	var rows []models.NonAgentUser
	err := config.DB.Find(&rows).Error
	return rows, err
}
