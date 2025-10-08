package models

import "time"

type Item struct {
	ID          int       `json:"id"`
	CompanyID   int       `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Type        string    `json:"type"` // "product" | "service"
	ItemPrice   float64   `json:"item_price"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (Item) TableName() string {
	return "items"
}
