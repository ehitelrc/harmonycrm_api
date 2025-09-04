package models

import "time"

type UserCompanyRole struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	CompanyID uint      `json:"company_id"`
	RoleID    uint      `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserCompanyRole) TableName() string { return "user_company_roles" }
