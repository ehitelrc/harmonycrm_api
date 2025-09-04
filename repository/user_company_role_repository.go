package repository

import (
	"errors"
	"harmony_api/config"
	"harmony_api/models"
)

type UserCompanyRoleRepository struct{}

// Estructuras de retorno
type UserPermissionRow struct {
	UserID         uint   `json:"user_id"`
	FullName       string `json:"full_name"`
	PermissionID   uint   `json:"permission_id"`
	PermissionCode string `json:"permission_code"`
	PermissionDesc string `json:"permission_description"`
	HasPermission  bool   `json:"has_permission"`
}

func NewUserCompanyRoleRepository() *UserCompanyRoleRepository { return &UserCompanyRoleRepository{} }

// Listados
func (r *UserCompanyRoleRepository) GetAll() ([]models.UserCompanyRole, error) {
	var rows []models.UserCompanyRole
	err := config.DB.Find(&rows).Error
	return rows, err
}
func (r *UserCompanyRoleRepository) GetByUser(userID uint) ([]models.UserCompanyRole, error) {
	var rows []models.UserCompanyRole
	err := config.DB.Where("user_id = ?", userID).Find(&rows).Error
	return rows, err
}
func (r *UserCompanyRoleRepository) GetByCompany(companyID uint) ([]models.UserCompanyRole, error) {
	var rows []models.UserCompanyRole
	err := config.DB.Where("company_id = ?", companyID).Find(&rows).Error
	return rows, err
}
func (r *UserCompanyRoleRepository) GetByRole(roleID uint) ([]models.UserCompanyRole, error) {
	var rows []models.UserCompanyRole
	err := config.DB.Where("role_id = ?", roleID).Find(&rows).Error
	return rows, err
}
func (r *UserCompanyRoleRepository) GetByID(id uint) (*models.UserCompanyRole, error) {
	var row models.UserCompanyRole
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Crear evitando duplicado lógico (user_id, company_id, role_id)
func (r *UserCompanyRoleRepository) Create(m *models.UserCompanyRole) error {
	var exists int64
	config.DB.Model(&models.UserCompanyRole{}).
		Where("user_id = ? AND company_id = ? AND role_id = ?", m.UserID, m.CompanyID, m.RoleID).
		Count(&exists)
	if exists > 0 {
		return errors.New("la combinación user_id + company_id + role_id ya existe")
	}
	return config.DB.Create(m).Error
}

// Update: recibe objeto completo con id
func (r *UserCompanyRoleRepository) Update(m *models.UserCompanyRole) error {
	// chequeo opcional de duplicado si cambian llaves
	var exists int64
	config.DB.Model(&models.UserCompanyRole{}).
		Where("user_id = ? AND company_id = ? AND role_id = ? AND id <> ?", m.UserID, m.CompanyID, m.RoleID, m.ID).
		Count(&exists)
	if exists > 0 {
		return errors.New("la combinación user_id + company_id + role_id ya existe en otro registro")
	}
	return config.DB.Save(m).Error
}

func (r *UserCompanyRoleRepository) Delete(id uint) error {
	return config.DB.Delete(&models.UserCompanyRole{}, id).Error
}

// 1) Usuarios y permisos por compañía y rol (matriz completa con has_permission)
func (r *UserCompanyRoleRepository) GetUsersAndPermissionsByCompanyRole(companyID, roleID uint) ([]UserPermissionRow, error) {
	rows := []UserPermissionRow{}
	err := config.DB.Raw(`
		SELECT
			u.id AS user_id,
			COALESCE(u.full_name, u.email) AS full_name,
			p.id AS permission_id,
			p.code AS permission_code,
			p.description AS permission_description,
			CASE WHEN rp.permission_id IS NOT NULL THEN TRUE ELSE FALSE END AS has_permission
		FROM user_company_roles ucr
		JOIN users u ON u.id = ucr.user_id
		CROSS JOIN permissions p
		LEFT JOIN role_permissions rp
			ON rp.role_id = ucr.role_id
		   AND rp.permission_id = p.id
		WHERE ucr.company_id = ?
		  AND ucr.role_id = ?
		ORDER BY u.id, p.code
	`, companyID, roleID).Scan(&rows).Error
	return rows, err
}

// 2) Permisos efectivos por compañía y usuario (unión de todos los roles del usuario en esa company)
func (r *UserCompanyRoleRepository) GetPermissionsByCompanyUser(companyID, userID uint) ([]models.Permission, error) {
	perms := []models.Permission{}
	err := config.DB.Raw(`
		SELECT DISTINCT p.id, p.code, p.description
		FROM user_company_roles ucr
		JOIN role_permissions rp ON rp.role_id = ucr.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ucr.company_id = ?
		  AND ucr.user_id = ?
		ORDER BY p.code
	`, companyID, userID).Scan(&perms).Error
	return perms, err
}
