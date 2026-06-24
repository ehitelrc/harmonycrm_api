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

	// Meta Integration fields
	MetaTemplateID *string `gorm:"column:meta_template_id" json:"meta_template_id,omitempty"`
	MetaCategory   *string `gorm:"column:meta_category" json:"meta_category,omitempty"`
	ApprovalStatus string  `gorm:"column:approval_status;default:local" json:"approval_status"`
	RejectionReason *string `gorm:"column:rejection_reason" json:"rejection_reason,omitempty"`
	BodyContent    *string `gorm:"column:body_content" json:"body_content,omitempty"`
	HeaderContent  *string `gorm:"column:header_content" json:"header_content,omitempty"`
	FooterContent  *string `gorm:"column:footer_content" json:"footer_content,omitempty"`
	ButtonsJSON    *string `gorm:"column:buttons_json" json:"buttons_json,omitempty"`

	// Calculated field — populated by GetAll, not stored in DB
	LinkedCount int `gorm:"column:linked_count;->" json:"linked_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides default table name
func (MessageTemplate) TableName() string {
	return "message_templates"
}
