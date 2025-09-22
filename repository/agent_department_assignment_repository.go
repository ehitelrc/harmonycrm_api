package repository

import (
	"harmony_api/config"
	"harmony_api/models"

	"gorm.io/gorm"
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

func (r *AgentDepartmentAssignmentRepository) GetByCompany(companyID uint) ([]models.VwAgentDepartmentAssignment, error) {
	var rows []models.VwAgentDepartmentAssignment
	err := config.DB.Where("company_id = ?", companyID).Find(&rows).Error
	return rows, err
}

func (r *AgentDepartmentAssignmentRepository) GetByCompanyAndAgent(companyID, agentID uint) ([]models.VwAgentDepartmentAssignment, error) {
	var rows []models.VwAgentDepartmentAssignment
	err := config.DB.Where("company_id = ? AND user_id = ?", companyID, agentID).Find(&rows).Error
	return rows, err
}

// Transactional delete all assignments for agent in company and create new ones
func (r *AgentDepartmentAssignmentRepository) SetAgentDepartments(companyID, agentID uint, assignments []models.VwAgentDepartmentAssignment) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Delete existing assignments for agent in company
		if err := tx.Where("agent_id = ? AND department_id IN (SELECT id FROM departments WHERE company_id = ?)", agentID, companyID).
			Delete(&models.AgentDepartmentAssignment{}).Error; err != nil {
			return err
		}

		// Create new assignments
		for _, a := range assignments {
			if a.DepartmentAssigned {
				newAssign := models.AgentDepartmentAssignment{
					AgentID:      agentID,
					DepartmentID: uint(a.DepartmentID),
				}
				if err := tx.Create(&newAssign).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *AgentDepartmentAssignmentRepository) GetAgentsByDepartment(companyID, departmentID uint) ([]models.AgentDepartmentInformation, error) {
	var rows []models.AgentDepartmentInformation

	err := config.DB.Where("company_id = ? AND department_id = ?", companyID, departmentID).Find(&rows).Error
	return rows, err
}
