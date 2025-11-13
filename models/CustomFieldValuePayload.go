package models

type CustomFieldValuePayload struct {
	EntityID   uint        `json:"entity_id"`
	EntityName string      `json:"entity_name"`
	FieldKey   string      `json:"field_key"`
	FieldValue interface{} `json:"field_value"`
}
