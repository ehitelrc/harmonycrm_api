// models/permission_view.go
package models

type UserEffectivePermission struct {
	UserID         int    `json:"user_id"`
	UserEmail      string `json:"user_email"`
	CompanyID      int    `json:"company_id"`
	RoleID         int    `json:"role_id"`
	RoleName       string `json:"role_name"`
	PermissionID   int    `json:"permission_id"`
	PermissionCode string `json:"permission_code"`
	PermissionDesc string `json:"permission_description"`
}

func (UserEffectivePermission) TableName() string {
	return "vw_user_effective_permissions"
}
