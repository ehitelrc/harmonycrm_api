package models

import "time"

type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Email           string    `json:"email" gorm:"unique;not null"`
	FullName        string    `json:"full_name"`
	Phone           string    `json:"phone"`
	ProfileImageURL string    `json:"profile_image_url"`
	PasswordHash    string    `json:"password_hash"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
