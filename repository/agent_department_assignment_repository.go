package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type AgentDepartmentAssignmentRepository struct{}

func NewAgentDepartmentAssignmentRepository() *AgentDepartmentAssignmentRepository {
	return &AgentDepartmentAssignmentRepository{}
}

func (r *AgentDepartmentAssignmentRepository) GetByID(id uint) (*models.AgentDepartmentAssignment, error) {
	var row models.AgentDepartmentAssignment
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *AgentDepartmentAssignmentRepository) GetByAgent(agentID uint) ([]models.AgentDepartmentAssignment, error) {
	var rows []models.AgentDepartmentAssignment
	err := config.DB.Where("agent_id = ?", agentID).Find(&rows).Error
	return rows, err
}

func (r *AgentDepartmentAssignmentRepository) GetByDepartment(deptID uint) ([]models.AgentDepartmentAssignment, error) {
	var rows []models.AgentDepartmentAssignment
	err := config.DB.Where("department_id = ?", deptID).Find(&rows).Error
	return rows, err
}

func (r *AgentDepartmentAssignmentRepository) Create(m *models.AgentDepartmentAssignment) error {
	return config.DB.Create(m).Error
}

// Update recibe el objeto completo (incluye id)
func (r *AgentDepartmentAssignmentRepository) Update(m *models.AgentDepartmentAssignment) error {
	return config.DB.Save(m).Error
}

func (r *AgentDepartmentAssignmentRepository) Delete(user_id uint) error {
	return config.DB.Delete(&models.AgentDepartmentAssignment{}, user_id).Error
}
