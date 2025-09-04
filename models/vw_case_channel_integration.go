package models

import "time"

// A) Modelo completo: para consultas internas (incluye access_token).
type VWCaseChannelIntegration struct {
	CaseID    int        `db:"case_id"                  json:"case_id"                   gorm:"column:case_id"`
	CompanyID *int       `db:"company_id"               json:"company_id,omitempty"      gorm:"column:company_id"`
	ChannelID *int64     `db:"channel_id"               json:"channel_id,omitempty"      gorm:"column:channel_id"`
	SenderID  *string    `db:"sender_id"                json:"sender_id,omitempty"       gorm:"column:sender_id"`
	Status    *string    `db:"status"                   json:"status,omitempty"          gorm:"column:status"`
	StartedAt *time.Time `db:"started_at"               json:"started_at,omitempty"      gorm:"column:started_at"`
	ClosedAt  *time.Time `db:"closed_at"                json:"closed_at,omitempty"       gorm:"column:closed_at"`

	ChannelIntegrationID int        `db:"channel_integration_id"   json:"channel_integration_id"    gorm:"column:channel_integration_id"`
	ChannelName          *string    `db:"channel_name"             json:"channel_name,omitempty"    gorm:"column:channel_name"`
	ChannelCode          string     `db:"channel_code"            json:"channel_code,omitempty"   gorm:"column:channel_code"`
	WebhookURL           string     `db:"webhook_url"              json:"webhook_url"               gorm:"column:webhook_url"`
	AccessToken          *string    `db:"access_token"             json:"access_token,omitempty"    gorm:"column:access_token"` // ⚠️ Sensible
	AppIdentifier        *string    `db:"app_identifier"           json:"app_identifier,omitempty"  gorm:"column:app_identifier"`
	IntegrationIsActive  bool       `db:"integration_is_active"    json:"integration_is_active"     gorm:"column:integration_is_active"`
	IntegrationUpdatedAt *time.Time `db:"integration_updated_at"   json:"integration_updated_at,omitempty" gorm:"column:integration_updated_at"`
}

// (Opcional) Para GORM: fija el nombre de la vista.
func (VWCaseChannelIntegration) TableName() string { return "vw_case_channel_integration" }

// B) DTO seguro: úsalo para responder por API (no expone access_token).
type CaseChannelIntegrationDTO struct {
	CaseID    int        `json:"case_id"`
	CompanyID *int       `json:"company_id,omitempty"`
	ChannelID *int64     `json:"channel_id,omitempty"`
	SenderID  *string    `json:"sender_id,omitempty"`
	Status    *string    `json:"status,omitempty"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`

	ChannelIntegrationID int        `json:"channel_integration_id"`
	ChannelName          *string    `json:"channel_name,omitempty"`
	WebhookURL           string     `json:"webhook_url"`
	AppIdentifier        *string    `json:"app_identifier,omitempty"`
	IntegrationIsActive  bool       `json:"integration_is_active"`
	IntegrationUpdatedAt *time.Time `json:"integration_updated_at,omitempty"`
}
