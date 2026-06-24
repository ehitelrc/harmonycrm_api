package models

type Channel struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Code        string  `gorm:"unique;not null" json:"code"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description"`
	MetaWabaID  *string `gorm:"column:meta_waba_id;type:text" json:"meta_waba_id,omitempty"`
}
