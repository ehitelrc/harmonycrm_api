package models

type VWCaseGeneralInformation struct {
	CaseID             uint    `json:"case_id" gorm:"column:case_id"`
	CompanyID          uint    `json:"company_id" gorm:"column:company_id"`
	ClientID           *uint   `json:"client_id,omitempty" gorm:"column:client_id"`
	DepartmentID       *uint   `json:"department_id,omitempty" gorm:"column:department_id"`
	AgentID            *uint   `json:"agent_id,omitempty" gorm:"column:agent_id"`
	FunnelID           *uint   `json:"funnel_id,omitempty" gorm:"column:funnel_id"`
	Status             string  `json:"status" gorm:"column:status"`
	ChannelID          *uint   `json:"channel_id,omitempty" gorm:"column:channel_id"`
	SenderID           *uint   `json:"sender_id,omitempty" gorm:"column:sender_id"`
	CurrentStageID     *uint   `json:"current_stage_id,omitempty" gorm:"column:current_stage_id"`
	CurrentStageName   *string `json:"current_stage_name,omitempty" gorm:"column:current_stage_name"`
	LastChangedByLabel *string `json:"last_changed_by_label,omitempty" gorm:"column:last_changed_by_label"`
	Action             *string `json:"action,omitempty" gorm:"column:action"`
	ClientName         *string `json:"client_name,omitempty" gorm:"column:client_name"`
	Email              *string `json:"email,omitempty" gorm:"column:email"`
	CampaignName       *string `json:"campaign_name,omitempty" gorm:"column:campaign_name"`
	ChannelName        *string `json:"channel_name,omitempty" gorm:"column:channel_name"`
	ChannelCode        *string `json:"channel_code,omitempty" gorm:"column:channel_code"`
	DepartmentName     *string `json:"department_name,omitempty" gorm:"column:department_name"`
	CampaignID         *uint   `json:"campaign_id,omitempty" gorm:"column:campaign_id"`
	AgentName          *string `json:"agent_name,omitempty" gorm:"column:agent_name"`
	ColorHex           *string `json:"color_hex,omitempty" gorm:"column:color_hex"`
}

func (VWCaseGeneralInformation) TableName() string {
	return "vw_case_general_information"
}
