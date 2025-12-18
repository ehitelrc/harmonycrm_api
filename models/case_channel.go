package models

import "time"

type CaseWithChannel struct {
	CaseID             int        `json:"case_id"               gorm:"column:case_id"`
	ClientID           *int       `json:"client_id,omitempty"   gorm:"column:client_id"`
	CampaignID         *int       `json:"campaign_id,omitempty" gorm:"column:campaign_id"`
	ClientName         *string    `json:"client_name" gorm:"column:client_name"`
	CompanyID          *int       `json:"company_id,omitempty"  gorm:"column:company_id"`
	DepartmentID       *int       `json:"department_id,omitempty" gorm:"column:department_id"`
	AgentID            *int       `json:"agent_id,omitempty"    gorm:"column:agent_id"`
	FunnelID           *int       `json:"funnel_id,omitempty"   gorm:"column:funnel_id"`
	FunnelStage        *string    `json:"funnel_stage,omitempty" gorm:"column:funnel_stage"`
	Status             *string    `json:"status,omitempty"      gorm:"column:status"`
	ChannelID          *int       `json:"channel_id,omitempty"  gorm:"column:channel_id"` // en cases es TEXT
	ChannelCode        *string    `json:"channel_code,omitempty" gorm:"column:channel_code"`
	ChannelName        *string    `json:"channel_name,omitempty" gorm:"column:channel_name"`
	ChannelDescription *string    `json:"channel_description,omitempty" gorm:"column:channel_description"`
	StartedAt          *time.Time `json:"started_at,omitempty"  gorm:"column:started_at"`
	ClosedAt           *time.Time `json:"closed_at,omitempty"   gorm:"column:closed_at"`
	CreatedAt          *time.Time `json:"created_at,omitempty"  gorm:"column:created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"  gorm:"column:updated_at"`
	SenderID           *string    `json:"sender_id,omitempty"   gorm:"column:sender_id"`

	LastMessageID         *uint      `json:"last_message_id" gorm:"column:last_message_id"`
	LastMessageSenderType *string    `json:"last_message_sender_type" gorm:"column:last_message_sender_type"` // 'client'|'agent'
	LastMessageType       *string    `json:"last_message_type" gorm:"column:last_message_type"`               // 'text'|'image'|'audio'|'file'
	LastMessageText       *string    `json:"last_message_text" gorm:"column:last_message_text"`
	LastMessageFileURL    *string    `json:"last_message_file_url" gorm:"column:last_message_file_url"`
	LastMessageMimeType   *string    `json:"last_message_mime_type" gorm:"column:last_message_mime_type"`
	LastMessageAt         *time.Time `json:"last_message_at" gorm:"column:last_message_at"`

	LastMessagePreview   *string `json:"last_message_preview" gorm:"column:last_message_preview"`
	LastMessageIsMedia   *bool   `json:"last_message_is_media" gorm:"column:last_message_is_media"`
	IsNonCommercial      *bool   `json:"is_non_commercial,omitempty" gorm:"column:is_non_commercial"`
	ManualStartingLead   *bool   `json:"manual_starting_lead,omitempty" gorm:"column:manual_starting_lead"`
	ClientMessages       *int    `json:"client_messages,omitempty" gorm:"column:client_messages"`
	ChannelIntegrationID uint    `json:"channel_integration_id" gorm:"column:channel_integration_id"`
	IntegrationName      string  `json:"integration_name" gorm:"column:integration_name"`

	// agent_assigned true si el caso tiene un agente asignado
	AgentAssigned bool `json:"agent_assigned" gorm:"column:agent_assigned"`

	//agent_full_name es el nombre completo del agente asignado
	AgentFullName *string `json:"agent_full_name,omitempty" gorm:"column:agent_full_name"`

	UnreadCount *int `json:"unread_count,omitempty" gorm:"column:unread_count"`
}

func (CaseWithChannel) TableName() string {
	return "vw_cases_with_channels" // nombre de la VISTA
}

type CaseWithChannelMV struct {
	CaseID             uint `gorm:"column:case_id"`
	ClientID           *uint
	ClientName         *string
	CampaignID         *uint
	CompanyID          uint
	DepartmentID       uint
	AgentID            *uint
	AgentAssigned      bool
	AgentFullName      *string
	FunnelID           *uint
	FunnelStage        *int
	Status             string
	ChannelID          uint
	ChannelCode        string
	ChannelName        string
	IntegrationName    *string
	StartedAt          *time.Time
	ClosedAt           *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	SenderID           string
	LastMessagePreview *string
	LastMessageAt      *time.Time
	ClientMessages     int64
	UnreadCount        int64
}

func (CaseWithChannelMV) TableName() string {
	return "mv_cases_with_channels"
}
