package models

type Setting struct {
	ID           int64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ValueCode    string   `json:"value_code" gorm:"type:varchar(200);not null"`
	TextValue    *string  `json:"text_value" gorm:"type:text"`
	IntegerValue *int64   `json:"integer_value"`
	NumberValue  *float64 `json:"number_value" gorm:"type:numeric(18,3)"`
	BoolValue    *bool    `json:"bool_value"`
	IsActive     bool     `json:"is_active" gorm:"not null"`
}

func (Setting) TableName() string {
	return "settings"
}
