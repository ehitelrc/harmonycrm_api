package repository

import (
	"errors"
	"harmony_api/config"
	"harmony_api/models"

	"gorm.io/gorm"
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

func (r *AgentRepository) CreateUnifiedAgent(email, fullName, phone, passwordHash string, companyID, roleID uint, departmentIDs []uint) (*models.User, error) {
	var user models.User
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Validar si ya existe el usuario
		var count int64
		if err := tx.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("el correo electrónico ya está registrado")
		}

		// 2. Crear el Usuario
		user = models.User{
			Email:        email,
			FullName:     fullName,
			Phone:        phone,
			PasswordHash: passwordHash,
			IsActive:     true,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 3. Asignar Rol en la Compañía
		roleAssign := models.UserCompanyRole{
			UserID:    user.ID,
			CompanyID: companyID,
			RoleID:    roleID,
		}
		if err := tx.Create(&roleAssign).Error; err != nil {
			return err
		}

		// 4. Registrar como Agente
		agent := models.Agent{
			UserID: user.ID,
		}
		if err := tx.Create(&agent).Error; err != nil {
			return err
		}

		// 5. Asignar Departamentos
		for _, deptID := range departmentIDs {
			assignment := models.AgentDepartmentAssignment{
				AgentID:      user.ID,
				DepartmentID: deptID,
			}
			if err := tx.Create(&assignment).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}
