package models

import "time"

type CustomFieldValue struct {
	ID         uint       `json:"id"`
	FieldID    uint       `json:"field_id"`
	EntityName string     `json:"entity_name"`
	EntityID   uint       `json:"entity_id"`
	ValueText  *string    `json:"value_text,omitempty"`
	ValueInt   *int       `json:"value_integer,omitempty"`
	ValueDec   *float64   `json:"value_decimal,omitempty"`
	ValueBool  *bool      `json:"value_boolean,omitempty"`
	ValueDate  *time.Time `json:"value_date,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
