package models

type ChannelWhatsAppTemplate struct {
	ID                 int64  `gorm:"primaryKey;column:id" json:"id"`
	DepartmentID       int64  `gorm:"column:department_id;not null" json:"department_id"`
	TemplateName       string `gorm:"column:template_name;size:50;not null" json:"template_name"`
	Language           string `gorm:"column:language;size:10;not null" json:"language"`
	Active             bool   `gorm:"column:active;default:true;not null" json:"active"`
	TemplateUrlWebhook string `gorm:"column:template_url_webhook" json:"template_url_webhook"`
}

// TableName especifica el nombre de la tabla en la BD
func (ChannelWhatsAppTemplate) TableName() string {
	return "channel_whatsapp_template"
}

type VwChannelWhatsAppTemplate struct {
	ID                 int64  `gorm:"primaryKey;column:id" json:"id"`
	TemplateName       string `gorm:"column:template_name;size:50;not null" json:"template_name"`
	Language           string `gorm:"column:language;size:10;not null" json:"language"`
	Active             bool   `gorm:"column:active;default:true;not null" json:"active"`
	TemplateUrlWebhook string `gorm:"column:template_url_webhook" json:"template_url_webhook"`

	CompanyID    int64 `gorm:"column:company_id;not null" json:"company_id"`
	ChannelID    int64 `gorm:"column:channel_id;not null" json:"channel_id"`
	DepartmentID int64 `gorm:"column:department_id;not null" json:"department_id"`
}

// TableName especifica el nombre de la tabla en la BD
func (VwChannelWhatsAppTemplate) TableName() string {
	return "vw_channel_whatsapp_templates"
}
