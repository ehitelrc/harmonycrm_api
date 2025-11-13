package models

import "time"

type CustomFieldValue struct {
	ID           uint       `json:"id"`
	FieldID      uint       `json:"field_id"`
	EntityName   string     `json:"entity_name"`
	EntityID     uint       `json:"entity_id"`
	ValueText    *string    `json:"value_text,omitempty"`
	ValueInteger *int       `json:"value_integer,omitempty"`
	ValueDecimal *float64   `json:"value_decimal,omitempty"`
	ValueBoolean *bool      `json:"value_boolean,omitempty"`
	ValueDate    *time.Time `json:"value_date,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (CustomFieldValue) TableName() string {
	return "custom_field_values"
}
