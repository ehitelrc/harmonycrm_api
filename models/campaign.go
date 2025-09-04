// models/campaign.go
package models

type Campaign struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	CompanyID   uint   `json:"company_id"`
	Name        string `json:"name" gorm:"not null"`
	StartDate   *Date  `json:"start_date" gorm:"type:date"`
	EndDate     *Date  `json:"end_date"   gorm:"type:date"`
	Description string `json:"description"`
	FunnelID    *uint  `json:"funnel_id"` // campaña puede no tener funnel
	IsActive    bool   `json:"is_active" gorm:"default:true"`
}

func (Campaign) TableName() string { return "campaigns" }
