package models

import "time"

type Client struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ExternalID string    `json:"external_id"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Client) TableName() string { return "clients" }
