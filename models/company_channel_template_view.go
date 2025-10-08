package models

type CompanyChannelTemplateView struct {
	CompanyID            int64   `json:"company_id" gorm:"column:company_id"`
	ChannelIntegrationID int64   `json:"channel_integration_id" gorm:"column:channel_integration_id"`
	WebhookURL           string  `json:"webhook_url" gorm:"column:webhook_url"`
	AccessToken          *string `json:"access_token,omitempty" gorm:"column:access_token"`
	AppIdentifier        *string `json:"app_identifier,omitempty" gorm:"column:app_identifier"`
	IntegrationActive    bool    `json:"integration_active" gorm:"column:integration_active"`
	IntegrationCreatedAt string  `json:"integration_created_at" gorm:"column:integration_created_at"`
	IntegrationUpdatedAt string  `json:"integration_updated_at" gorm:"column:integration_updated_at"`

	ChannelID          int64   `json:"channel_id" gorm:"column:channel_id"`
	ChannelCode        string  `json:"channel_code" gorm:"column:channel_code"`
	ChannelName        string  `json:"channel_name" gorm:"column:channel_name"`
	ChannelDescription *string `json:"channel_description,omitempty" gorm:"column:channel_description"`

	TemplateID     *int64  `json:"template_id,omitempty" gorm:"column:template_id"`
	TemplateName   *string `json:"template_name,omitempty" gorm:"column:template_name"`
	Language       *string `json:"language,omitempty" gorm:"column:language"`
	TemplateActive *bool   `json:"template_active,omitempty" gorm:"column:template_active"`

	TemplateUrlWebhook string `json:"template_url_webhook" gorm:"column:template_url_webhook"`
}

func (CompanyChannelTemplateView) TableName() string {
	return "public.vw_company_channel_templates"
}
