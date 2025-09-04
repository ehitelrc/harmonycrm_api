package models

import "time"

type VWChannel struct {
	IntegrationID      uint      `json:"integration_id" gorm:"column:integration_id"`
	CompanyID          uint      `json:"company_id" gorm:"column:company_id"`
	ChannelID          string    `json:"channel_id" gorm:"column:channel_id"`
	ChannelCode        string    `json:"channel_code" gorm:"column:channel_code"`
	ChannelName        string    `json:"channel_name" gorm:"column:channel_name"`
	ChannelDescription string    `json:"channel_description" gorm:"column:channel_description"`
	WebhookURL         string    `json:"webhook_url" gorm:"column:webhook_url"`
	AccessToken        string    `json:"access_token" gorm:"column:access_token"`
	AppIdentifier      string    `json:"app_identifier" gorm:"column:app_identifier"`
	IsActive           bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName sets the table name for VWChannel to the view name
func (VWChannel) TableName() string {
	return "vw_channels"
}
