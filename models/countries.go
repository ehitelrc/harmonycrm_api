package models

import (
	"encoding/json"
	"time"
)

type Country struct {
	ID           uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string        `gorm:"type:text;not null" json:"name"`
	ISOCode      string        `gorm:"type:char(2);uniqueIndex;not null" json:"iso_code"`
	PhoneCode    string        `gorm:"type:text" json:"phone_code"`
	CurrencyCode string        `gorm:"type:char(3)" json:"currency_code"`
	Provinces    []Province    `gorm:"foreignKey:CountryCode;references:ISOCode;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"provinces,omitempty"`
	CreatedAt    *TimeOnlyJSON `json:"created_at,omitempty"`
	UpdatedAt    *TimeOnlyJSON `json:"updated_at,omitempty"`
}

func (Country) TableName() string { return "countries" }

// Province representa una provincia

type Province struct {
	ID             uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	CountryCode    string        `gorm:"type:char(2);not null;index" json:"country_code"`
	Code           string        `gorm:"type:text;not null;index" json:"code"`
	Name           string        `gorm:"type:text;not null" json:"name"`
	ProvinceNumber int64         `gorm:"not null" json:"province_number"`
	Country        Country       `gorm:"foreignKey:CountryCode;references:ISOCode;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"country,omitempty"`
	Cantons        []Canton      `gorm:"foreignKey:ProvinceCode;references:Code;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"cantons,omitempty"`
	CreatedAt      *TimeOnlyJSON `json:"created_at,omitempty"`
	UpdatedAt      *TimeOnlyJSON `json:"updated_at,omitempty"`
}

func (Province) TableName() string { return "provinces" }

// Canton representa un cantón
type Canton struct {
	ID           uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	CountryCode  string        `gorm:"type:char(2);not null;index" json:"country_code"`
	ProvinceCode string        `gorm:"type:text;not null;index" json:"province_code"`
	Code         string        `gorm:"type:text;not null;uniqueIndex" json:"code"`
	Name         string        `gorm:"type:text;not null" json:"name"`
	CantonNumber int64         `gorm:"not null" json:"canton_number"`
	Province     Province      `gorm:"foreignKey:ProvinceCode;references:Code;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"province,omitempty"`
	Districts    []District    `gorm:"foreignKey:CantonCode;references:Code;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"districts,omitempty"`
	CreatedAt    *TimeOnlyJSON `json:"created_at,omitempty"`
	UpdatedAt    *TimeOnlyJSON `json:"updated_at,omitempty"`
}

func (Canton) TableName() string { return "cantons" }

// District representa un distrito
type District struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	CountryCode string        `gorm:"type:char(2);not null;index" json:"country_code"`
	CantonCode  string        `gorm:"type:text;not null;index" json:"canton_code"`
	Code        string        `gorm:"type:text;not null;uniqueIndex" json:"code"`
	Name        string        `gorm:"type:text;not null" json:"name"`
	Latitude    *float64      `gorm:"type:numeric(10,7)" json:"latitude,omitempty"`
	Longitude   *float64      `gorm:"type:numeric(10,7)" json:"longitude,omitempty"`
	PostalCode  *string       `gorm:"type:text" json:"postal_code,omitempty"`
	Canton      Canton        `gorm:"foreignKey:CantonCode;references:Code;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"canton,omitempty"`
	CreatedAt   *TimeOnlyJSON `json:"created_at,omitempty"`
	UpdatedAt   *TimeOnlyJSON `json:"updated_at,omitempty"`
}

func (District) TableName() string { return "districts" }

type TimeOnlyJSON time.Time

func (t *TimeOnlyJSON) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	stamp := time.Time(*t).Format(`"2006-01-02 15:04:05"`)
	return []byte(stamp), nil
}

func (t *TimeOnlyJSON) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsed, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}
	*t = TimeOnlyJSON(parsed)
	return nil
}
