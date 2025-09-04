package models

import "time"

// Mapea filas de la vista public.vw_campaigns_with_funnel
type CampaignWithFunnel struct {
	CampaignID   uint       `json:"campaign_id"  gorm:"column:campaign_id"`
	CompanyID    *uint      `json:"company_id"   gorm:"column:company_id"` // NULLable
	CampaignName string     `json:"campaign_name" gorm:"column:campaign_name"`
	StartDate    *time.Time `json:"start_date"   gorm:"column:start_date"`  // DATE, NULLable
	EndDate      *time.Time `json:"end_date"     gorm:"column:end_date"`    // DATE, NULLable
	Description  *string    `json:"description"  gorm:"column:description"` // NULLable
	IsActive     bool       `json:"is_active"    gorm:"column:is_active"`
	FunnelID     *uint      `json:"funnel_id"    gorm:"column:funnel_id"`   // NULLable
	FunnelName   *string    `json:"funnel_name"  gorm:"column:funnel_name"` // NULLable
}

// Indicar a GORM que la “tabla” es la vista.
func (CampaignWithFunnel) TableName() string {
	return "vw_campaigns_with_funnel"
}
