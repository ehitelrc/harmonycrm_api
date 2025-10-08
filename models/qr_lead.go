package models

import "time"

type QrLead struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID    int       `json:"company_id" gorm:"not null"`
	CampaignID   int       `json:"campaign_id" gorm:"not null"`
	DepartmentID *int      `json:"department_id"`
	UserID       *int      `json:"user_id"`
	ClientID     *uint     `json:"client_id"`
	ContactPhone string    `json:"contact_phone" gorm:"type:varchar(50);not null"`
	Status       string    `json:"status" gorm:"type:text;default:'pending'"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones opcionales
	Company    *Company    `json:"company" gorm:"foreignKey:CompanyID"`
	Campaign   *Campaign   `json:"campaign" gorm:"foreignKey:CampaignID"`
	Department *Department `json:"department" gorm:"foreignKey:DepartmentID"`
	User       *User       `json:"user" gorm:"foreignKey:UserID"`
	Client     *Client     `json:"client" gorm:"foreignKey:ClientID"`
}

func (QrLead) TableName() string {
	return "qr_leads"
}
