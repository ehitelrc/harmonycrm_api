package models

import "time"

type CustomFieldDefinition struct {
	ID         uint      `json:"id"`
	EntityName string    `json:"entity_name"`
	FieldKey   string    `json:"field_key"`
	Label      string    `json:"label"`
	FieldType  string    `json:"field_type"`
	IsRequired bool      `json:"is_required"`
	IsActive   bool      `json:"is_active"`
	SortOrder  int       `json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
