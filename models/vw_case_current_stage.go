package models

import "time"

// VWCaseCurrentStage mapea la vista public.vw_case_current_stage
type VWCaseCurrentStage struct {
	CaseID             int       `gorm:"column:case_id"              json:"case_id"`
	FunnelID           *int      `gorm:"column:funnel_id"            json:"funnel_id,omitempty"`
	FunnelName         *string   `gorm:"column:funnel_name"          json:"funnel_name,omitempty"`
	CurrentStageID     *int      `gorm:"column:current_stage_id"     json:"current_stage_id,omitempty"`
	CurrentStageName   *string   `gorm:"column:current_stage_name"   json:"current_stage_name,omitempty"`
	LastChangedAt      time.Time `gorm:"column:last_changed_at"      json:"last_changed_at"`
	LastChangedBy      int       `gorm:"column:last_changed_by"      json:"last_changed_by"`
	LastChangedByLabel *string   `gorm:"column:last_changed_by_label" json:"last_changed_by_label,omitempty"`
	Action             string    `gorm:"column:action"               json:"action"`
}

// TableName especifica el nombre de la vista en la BD.
func (VWCaseCurrentStage) TableName() string {
	return "vw_case_current_stage"
}
