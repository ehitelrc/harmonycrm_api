// models/case_funnel.go
package models

import "time"

type CaseFunnel struct {
	ID          uint      `json:"id"            gorm:"primaryKey;autoIncrement"`
	CaseID      int       `json:"case_id"       gorm:"not null;index"`
	FunnelID    int       `json:"funnel_id"     gorm:"not null;index"`
	FromStageID *int      `json:"from_stage_id" gorm:"index"`     // nullable
	ToStageID   *int      `json:"to_stage_id"   gorm:"index"`     // nullable (IMPORTANTE)
	Note        *string   `json:"note"          gorm:"type:text"` // nullable
	ChangedBy   int       `json:"changed_by"    gorm:"not null;index"`
	ChangedAt   time.Time `json:"changed_at"    gorm:"not null;default:now()"`
	Action      string    `json:"action"        gorm:"type:text;not null;default:'move'"`
}

func (CaseFunnel) TableName() string { return "case_funnel" }
