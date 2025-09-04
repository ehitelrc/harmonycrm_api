package models

type VWUserPermission struct {
	UserID                int     `gorm:"column:user_id" json:"user_id"`
	CompanyID             int     `gorm:"column:company_id" json:"company_id"`
	PermissionID          int     `gorm:"column:permission_id" json:"permission_id"`
	PermissionCode        string  `gorm:"column:permission_code" json:"permission_code"`
	PermissionDescription *string `gorm:"column:permission_description" json:"permission_description,omitempty"`
	HasPermission         bool    `gorm:"column:has_permission" json:"has_permission"`
}

func (VWUserPermission) TableName() string { return "vw_user_permissions" }
