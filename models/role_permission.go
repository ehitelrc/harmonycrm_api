package models

type RolePermission struct {
	RoleID       uint `json:"role_id" gorm:"primaryKey"`
	PermissionID uint `json:"permission_id" gorm:"primaryKey"`
}

func (RolePermission) TableName() string { return "role_permissions" }
