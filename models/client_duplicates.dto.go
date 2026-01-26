// models/client_duplicates.dto.go
package models

type ClientDTO struct {
	ID            uint    `json:"id"`
	ExternalID    *string `json:"external_id"`
	FullName      *string `json:"full_name"`
	Email         *string `json:"email"`
	Phone         *string `json:"phone"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	CountryID     *uint   `json:"country_id"`
	ProvinceID    *uint   `json:"province_id"`
	CantonID      *uint   `json:"canton_id"`
	DistrictID    *uint   `json:"district_id"`
	AddressDetail *string `json:"address_detail"`
	PostalCode    *string `json:"postal_code"`
	IsCitizen     bool    `json:"is_citizen"`
}

type DuplicatePhoneGroupDTO struct {
	Phone   string      `json:"phone"`
	Count   int         `json:"count"`
	Clients []ClientDTO `json:"clients"`
}
