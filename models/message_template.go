package models

import "time"

type MessageTemplate struct {
	ID uint `gorm:"primaryKey" json:"id"`

	ChannelID uint `gorm:"not null;index" json:"channel_id"`

	TemplateName string `gorm:"type:text;not null" json:"template_name"`

	LanguageCode string `gorm:"type:varchar(10);not null" json:"language_code"`

	Description string `gorm:"type:text" json:"description"`

	Category string `gorm:"type:text" json:"category"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	IsConversationStarter bool `gorm:"default:false" json:"is_conversation_starter"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides default table name
func (MessageTemplate) TableName() string {
	return "message_templates"
}
