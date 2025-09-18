package models

type UserRoleCompanyManage struct {
	UserID          uint   `json:"user_id"  `
	CompanyID       uint   `json:"company_id"  `
	RoleID          uint   `json:"role_id"  `
	UserRoleCompany string `json:"user_role_company"  `
	HasRole         bool   `json:"has_role"  `
}
