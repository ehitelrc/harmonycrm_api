package models

type CustomListDefinition struct {
	ID               uint   `json:"id" gorm:"primaryKey"`
	ListName         string `json:"list_name"`
	CodeLabel        string `json:"code_label"`
	DescriptionLabel string `json:"description_label"`
	EntityName       string `json:"entity_name"`
	ListLabel        string `json:"list_label"`
}

type CustomListValue struct {
	ID               uint   `json:"id" gorm:"primaryKey"`
	ListID           uint   `json:"list_id"`
	CodeValue        string `json:"code_value"`
	DescriptionValue string `json:"description_value"`
}

type CustomListEntityValue struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	EntityName string `json:"entity_name"`
	EntityID   uint   `json:"entity_id"`
	ListValue  uint   `json:"list_value"` // -> refiere a custom_list_values.id
	ListID     uint   `json:"list_id"`
}

func (CustomListEntityValue) TableName() string {
	return "custom_list_entity_value"
}

/*
 * DTO final que consume ANGULAR:
 */
type CustomListDTO struct {
	ListID           uint                 `json:"list_id"`
	ListName         string               `json:"list_name"`
	CodeLabel        string               `json:"code_label"`
	DescriptionLabel string               `json:"description_label"`
	Values           []CustomListValueDTO `json:"values"`
	SelectedValue    *uint                `json:"selected_value"`
}

type CustomListValueDTO struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}
