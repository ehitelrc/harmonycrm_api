package repository

import (
	"harmony_api/config"
	"harmony_api/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// Assign en lote de forma transaccional
func (r *RolePermissionRepository) AssignBatch(assignments []models.AssignRequest) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// roleID -> slice de permissionIDs que deben QUEDAR asignados
		keepByRole := make(map[uint][]uint)
		// set de roles vistos en el payload
		rolesSeen := make(map[uint]struct{})

		for _, a := range assignments {
			rolesSeen[a.RoleID] = struct{}{}

			// Si tu frontend envía true/false, solo guardamos los true.
			// Si tu frontend envía SOLO los asignados (true), también caen aquí.
			if a.AssignRequest {
				keepByRole[a.RoleID] = append(keepByRole[a.RoleID], a.PermissionID)
			}
		}

		// Recorremos cada rol afectado
		for roleID := range rolesSeen {
			idsToKeep := keepByRole[roleID]

			// 1) Borrar lo que NO esté en la lista de "keep".
			if len(idsToKeep) > 0 {
				if err := tx.
					Where("role_id = ? AND permission_id NOT IN ?", roleID, idsToKeep).
					Delete(&models.RolePermission{}).Error; err != nil {
					return err
				}
			} else {
				// Si no hay nada que conservar, se eliminan TODOS los permisos del rol
				if err := tx.
					Where("role_id = ?", roleID).
					Delete(&models.RolePermission{}).Error; err != nil {
					return err
				}
			}

			// 2) Insertar los que deben quedar, ignorando duplicados
			if len(idsToKeep) > 0 {
				batch := make([]models.RolePermission, 0, len(idsToKeep))
				for _, pid := range idsToKeep {
					batch = append(batch, models.RolePermission{
						RoleID:       roleID,
						PermissionID: pid,
					})
				}

				if err := tx.
					Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "role_id"}, {Name: "permission_id"}},
						DoNothing: true,
					}).
					Create(&batch).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Unassign elimina la asignación de un permiso a un rol

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

// GetRolePermissionsViewByRole obtiene los permisos de un rol con info adicional de la vista
func (r *RolePermissionRepository) GetViewByRole(roleID uint) ([]models.RolePermissionView, error) {
	var rows []models.RolePermissionView
	err := config.DB.Where("role_id = ?", roleID).Find(&rows).Error
	return rows, err
}
