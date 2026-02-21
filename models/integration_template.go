package models

import "time"

type IntegrationTemplate struct {
	ID uint `gorm:"primaryKey" json:"id"`

	IntegrationID uint `gorm:"not null;index" json:"integration_id"`

	TemplateID uint `gorm:"not null;index" json:"template_id"`

	CreatedAt time.Time `json:"created_at"`
}

// TableName overrides default table name
func (IntegrationTemplate) TableName() string {
	return "integration_templates"
}
