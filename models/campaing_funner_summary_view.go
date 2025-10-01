package models

// CampaignFunnelSummary representa el resumen del funnel por campaña.
type CampaignFunnelSummary struct {
	CampaignID   int    `json:"campaign_id" gorm:"column:campaign_id"`
	CampaignName string `json:"campaign_name" gorm:"column:campaign_name"`
	CompanyID    int    `json:"company_id" gorm:"column:company_id"`
	FunnelID     int    `json:"funnel_id" gorm:"column:funnel_id"`
	FunnelName   string `json:"funnel_name" gorm:"column:funnel_name"`
	StageID      int    `json:"stage_id" gorm:"column:stage_id"`
	StageName    string `json:"stage_name" gorm:"column:stage_name"`
	StageCode    string `json:"stage_code" gorm:"column:stage_code"`
	Position     int    `json:"position" gorm:"column:position"`
	ColorHex     string `json:"color_hex" gorm:"column:color_hex"`
	IsWon        bool   `json:"is_won" gorm:"column:is_won"`
	IsLost       bool   `json:"is_lost" gorm:"column:is_lost"`
	IsTerminal   bool   `json:"is_terminal" gorm:"column:is_terminal"`
	TotalCases   int    `json:"total_cases" gorm:"column:total_cases"`
}

// TableName sobreescribe el nombre de la vista en la BD.
func (CampaignFunnelSummary) TableName() string {
	return "vw_campaign_funnel_summary"
}
