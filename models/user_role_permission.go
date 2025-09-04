package models

type UserRolePermission struct {
	CompanyID             int    `json:"company_id"`
	UserID                int    `json:"user_id"`
	UserName              string `json:"user_name"`
	RoleID                int    `json:"role_id"`
	RoleName              string `json:"role_name"`
	IsAgent               bool   `json:"is_agent"`
	PermissionID          int    `json:"permission_id"`
	PermissionCode        string `json:"permission_code"`
	PermissionDescription string `json:"permission_description"`
	HasPermission         bool   `json:"has_permission"`
}

func (m *UserRolePermission) TableName() string {
	return "vw_user_role_permissions"
}
