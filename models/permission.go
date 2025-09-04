package models

type Permission struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"unique;not null"`
	Description string `json:"description"`
}

func (Permission) TableName() string { return "permissions" }
