package models

type VwDashboardCampaignPerCompany struct {
	CompanyID      int64   `json:"company_id" gorm:"column:company_id"`
	CompanyName    string  `json:"company_name" gorm:"column:company_name"`
	CampaignID     int64   `json:"campaign_id" gorm:"column:campaign_id"`
	CampaignName   string  `json:"campaign_name" gorm:"column:campaign_name"`
	IsActive       bool    `json:"is_active" gorm:"column:is_active"`
	TotalCases     int64   `json:"total_cases" gorm:"column:total_cases"`
	OpenCases      int64   `json:"open_cases" gorm:"column:open_cases"`
	ClosedCases    int64   `json:"closed_cases" gorm:"column:closed_cases"`
	WonCases       int64   `json:"won_cases" gorm:"column:won_cases"`
	LostCases      int64   `json:"lost_cases" gorm:"column:lost_cases"`
	ConversionRate float64 `json:"conversion_rate" gorm:"column:conversion_rate"`
}

func (VwDashboardCampaignPerCompany) TableName() string {
	return "vw_dashboard_campaigns_per_company"
}
