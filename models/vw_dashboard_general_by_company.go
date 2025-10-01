package models

type VwDashboardGeneralByCompany struct {
	CompanyID            int64   `json:"company_id" gorm:"column:company_id"`
	CompanyName          string  `json:"company_name" gorm:"column:company_name"`
	TotalActiveCampaigns int64   `json:"total_active_campaigns" gorm:"column:total_active_campaigns"`
	TotalCases           int64   `json:"total_cases" gorm:"column:total_cases"`
	ClosedCases          int64   `json:"closed_cases" gorm:"column:closed_cases"`
	WonCases             int64   `json:"won_cases" gorm:"column:won_cases"`
	ConversionRate       float64 `json:"conversion_rate" gorm:"column:conversion_rate"`
	OperatingAgents      int64   `json:"operating_agents" gorm:"column:operating_agents"`
}

func (VwDashboardGeneralByCompany) TableName() string {
	return "vw_dashboard_general_by_company"
}
