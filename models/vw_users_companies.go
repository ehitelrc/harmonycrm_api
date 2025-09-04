package models

// VWUsersCompanies representa la vista public.vw_users_companies
type VWUsersCompanies struct {
	UserID      int     `gorm:"column:user_id" json:"user_id"`
	Email       string  `gorm:"column:email" json:"email"`
	FullName    *string `gorm:"column:full_name" json:"full_name,omitempty"`
	Phone       *string `gorm:"column:phone" json:"phone,omitempty"`
	IsActive    *bool   `gorm:"column:is_active" json:"is_active,omitempty"`
	CompanyID   int     `gorm:"column:company_id" json:"company_id"`
	CompanyName string  `gorm:"column:company_name" json:"company_name"`
	Token       string  `gorm:"column:token" json:"token"`
}

func (VWUsersCompanies) TableName() string {
	return "vw_users_companies"
}
