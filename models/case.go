package models

import "time"

type Case struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ClientID     *uint      `gorm:"default:null" json:"client_id,omitempty"`
	CampaignID   uint       `gorm:"default:null" json:"campaign_id"`
	CompanyID    uint       `gorm:"default:null" json:"company_id"`
	DepartmentID uint       `gorm:"default:null" json:"department_id"`
	AgentID      uint       `gorm:"default:null" json:"agent_id"`
	FunnelID     uint       `gorm:"default:null" json:"funnel_id"`
	FunnelStage  string     `gorm:"size:100" json:"funnel_stage"`
	Status       string     `gorm:"type:enum('open','in_progress','closed','cancelled')" json:"status"`
	ChannelID    string     `gorm:"default:null" json:"channel_id"`
	SenderId     string     `gorm:"default:null" json:"sender_id"`
	StartedAt    time.Time  `json:"started_at"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
