package models

type ChannelWhatsAppTemplate struct {
	ID                 int64  `gorm:"primaryKey;column:id" json:"id"`
	ChannelIntegration int64  `gorm:"column:channel_integration;not null" json:"channel_integration"`
	TemplateName       string `gorm:"column:template_name;size:50;not null" json:"template_name"`
	Language           string `gorm:"column:language;size:10;not null" json:"language"`
	Active             bool   `gorm:"column:active;default:true;not null" json:"active"`
}

// TableName especifica el nombre de la tabla en la BD
func (ChannelWhatsAppTemplate) TableName() string {
	return "channel_whatsapp_template"
}
