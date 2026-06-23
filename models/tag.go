package models

import "time"

type Tag struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Color        string    `gorm:"size:50;not null" json:"color"`
	Icon         string    `gorm:"size:100;not null" json:"icon"`
	DepartmentID *uint     `gorm:"default:null" json:"department_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Tag) TableName() string { return "tags" }
