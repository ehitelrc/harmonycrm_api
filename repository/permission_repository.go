package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type PermissionRepository struct{}

func NewPermissionRepository() *PermissionRepository { return &PermissionRepository{} }

func (r *PermissionRepository) GetAll() ([]models.Permission, error) {
	var rows []models.Permission
	err := config.DB.Find(&rows).Error
	return rows, err
}

func (r *PermissionRepository) GetByID(id uint) (*models.Permission, error) {
	var row models.Permission
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *PermissionRepository) GetByUserID(userID uint, companyID uint) ([]models.VWUserPermission, error) {

	var rows []models.VWUserPermission
	err := config.DB.
		Raw(`SELECT * FROM vw_user_permissions WHERE user_id = ? AND company_id = ?`, userID, companyID).
		Scan(&rows).Error

	return rows, err
}

func (r *PermissionRepository) Create(m *models.Permission) error {
	return config.DB.Create(m).Error
}

// Update recibe el objeto completo (incluye id)
func (r *PermissionRepository) Update(m *models.Permission) error {
	return config.DB.Save(m).Error
}

func (r *PermissionRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Permission{}, id).Error
}
