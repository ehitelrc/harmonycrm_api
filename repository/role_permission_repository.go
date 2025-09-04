package repository

import (
	"harmony_api/config"
	"harmony_api/models"

	"gorm.io/gorm"
)

type RolePermissionRepository struct{}

func NewRolePermissionRepository() *RolePermissionRepository { return &RolePermissionRepository{} }

func (r *RolePermissionRepository) GetByRole(roleID uint) ([]models.RolePermission, error) {
	var rows []models.RolePermission
	err := config.DB.Where("role_id = ?", roleID).Find(&rows).Error
	return rows, err
}

func (r *RolePermissionRepository) GetByPermission(permissionID uint) ([]models.RolePermission, error) {
	var rows []models.RolePermission
	err := config.DB.Where("permission_id = ?", permissionID).Find(&rows).Error
	return rows, err
}

func (r *RolePermissionRepository) Assign(roleID, permissionID uint) error {
	rp := models.RolePermission{RoleID: roleID, PermissionID: permissionID}
	return config.DB.Create(&rp).Error
}

func (r *RolePermissionRepository) Unassign(roleID, permissionID uint) error {
	return config.DB.
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{}).Error
}

// ReplaceRolePermissions borra las asignaciones actuales y crea las nuevas en una sola transacción
func (r *RolePermissionRepository) ReplaceRolePermissions(roleID uint, permissionIDs []uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// eliminar existentes del rol
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		// insertar nuevas (si la lista está vacía, queda sin permisos)
		if len(permissionIDs) == 0 {
			return nil
		}
		batch := make([]models.RolePermission, 0, len(permissionIDs))
		for _, pid := range permissionIDs {
			batch = append(batch, models.RolePermission{RoleID: roleID, PermissionID: pid})
		}
		return tx.Create(&batch).Error
	})
}
