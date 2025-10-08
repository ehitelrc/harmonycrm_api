package models

import (
	"time"
)

type ClientCampaign struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID        *uint     `gorm:"column:client_id" json:"client_id,omitempty"`
	CampaignID      *uint     `gorm:"column:campaign_id" json:"campaign_id,omitempty"`
	FunnelID        *uint     `gorm:"column:funnel_id" json:"funnel_id,omitempty"`
	StageDetail     *string   `gorm:"column:stage_detail" json:"stage_detail,omitempty"`
	AssignedAgentID *uint     `gorm:"column:assigned_agent_id" json:"assigned_agent_id,omitempty"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relaciones (opcionales)
	Client   *Client   `gorm:"foreignKey:ClientID;references:ID" json:"client,omitempty"`
	Campaign *Campaign `gorm:"foreignKey:CampaignID;references:ID" json:"campaign,omitempty"`
	Funnel   *Funnel   `gorm:"foreignKey:FunnelID;references:ID" json:"funnel,omitempty"`
	Agent    *User     `gorm:"foreignKey:AssignedAgentID;references:ID" json:"agent,omitempty"`
}

func (ClientCampaign) TableName() string {
	return "client_campaigns"
}
