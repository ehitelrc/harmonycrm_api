package models

import "time"

type Client struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	ExternalID string `json:"external_id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`

	CountryID     uint   `json:"country_id"`
	ProvinceID    uint   `json:"province_id"`
	CantonID      uint   `json:"canton_id"`
	DistrictID    uint   `json:"district_id"`
	AddressDetail string `json:"address_detail"`
	PostalCode    string `json:"postal_code"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Client) TableName() string { return "clients" }
