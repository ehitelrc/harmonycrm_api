package models

import "time"

// VWChannelIntegration representa una fila de la vista public.vw_channel_integrations
type VWChannelIntegration struct {
	ChannelIntegrationID int64     `json:"channel_integration_id" gorm:"column:channel_integration_id"`
	CompanyID            int64     `json:"company_id"             gorm:"column:company_id"`
	ChannelID            int64     `json:"channel_id"             gorm:"column:channel_id"`
	ChannelName          string    `json:"channel_name"           gorm:"column:channel_name"`
	ChannelCode          string    `json:"channel_code"           gorm:"column:channel_code"`
	WebhookURL           string    `json:"webhook_url"            gorm:"column:webhook_url"`
	AccessToken          *string   `json:"access_token,omitempty" gorm:"column:access_token"`
	AppIdentifier        *string   `json:"app_identifier,omitempty" gorm:"column:app_identifier"`
	IsActive             bool      `json:"is_active"              gorm:"column:is_active"`
	CreatedAt            time.Time `json:"created_at"             gorm:"column:created_at"`
	UpdatedAt            time.Time `json:"updated_at"             gorm:"column:updated_at"`
}

// TableName asegura que GORM consulte la vista correcta
func (VWChannelIntegration) TableName() string {
	return "public.vw_channel_integrations"
}
