package models

type ChannelTemplateIntegration struct {
	ChannelID   uint   `gorm:"column:channel_id"       json:"channel_id"`
	ChannelCode string `gorm:"column:channel_code"     json:"channel_code"`
	ChannelName string `gorm:"column:channel_name"     json:"channel_name"`

	TemplateID   uint   `gorm:"column:template_id"   json:"template_id"`
	TemplateName string `gorm:"column:template_name" json:"template_name"`
	LanguageCode string `gorm:"column:language_code" json:"language_code"`

	IntegrationID   uint   `gorm:"column:integration_id"   json:"integration_id"`
	IntegrationName string `gorm:"column:integration_name" json:"integration_name"`
	CompanyID       uint   `gorm:"column:company_id"       json:"company_id"`
	DepartmentID    *uint  `gorm:"column:department_id"    json:"department_id"`

	IsLinked bool `gorm:"column:is_linked" json:"is_linked"`
}

func (ChannelTemplateIntegration) TableName() string {
	return "vw_channel_template_integrations"
}
