package models

import "time"

type Company struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Industry  string    `gorm:"size:255" json:"industry"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Company) TableName() string {
	return "companies"
}
