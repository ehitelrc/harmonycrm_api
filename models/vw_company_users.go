package models

type CompanyUserView struct {
	CompanyID   uint   `json:"company_id" gorm:"column:company_id"`
	CompanyName string `json:"company_name" gorm:"column:company_name"`
	UserID      uint   `json:"user_id" gorm:"column:user_id"`
	FullName    string `json:"full_name" gorm:"column:full_name"`
	Email       string `json:"email" gorm:"column:email"`
}

// Nombre real de la vista en la base de datos
func (CompanyUserView) TableName() string {
	return "vw_company_users"
}
