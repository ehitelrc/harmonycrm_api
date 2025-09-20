package models

// RolePermissionView representa la vista public.role_permissions_view
// Es de solo lectura (->) y NO tiene primary key porque es una vista.
type RolePermissionView struct {
	RoleID       uint    `json:"role_id" gorm:"column:role_id;->"`
	PermissionID uint    `json:"permission_id" gorm:"column:permission_id;->"`
	Code         string  `json:"code" gorm:"column:code;->"`
	Description  *string `json:"description,omitempty" gorm:"column:description;->"`
	Assigned     bool    `json:"assigned" gorm:"column:assigned;->"`
}

func (RolePermissionView) TableName() string { return "vw_role_permissions" }
