package models

import "time"

// Country representa un país
type Country struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:text;not null" json:"name"`
	ISOCode      string    `gorm:"type:char(2);unique" json:"iso_code"`
	PhoneCode    string    `gorm:"type:text" json:"phone_code"`
	CurrencyCode string    `gorm:"type:char(3)" json:"currency_code"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Country) TableName() string { return "countries" }

// Province representa una provincia
type Province struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:text;not null" json:"name"`
	CountryID uint      `gorm:"not null" json:"country_id"`
	Country   Country   `gorm:"foreignKey:CountryID;constraint:OnDelete:CASCADE;" json:"country,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Province) TableName() string { return "provinces" }

// Canton representa un cantón
type Canton struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:text;not null" json:"name"`
	ProvinceID uint      `gorm:"not null" json:"province_id"`
	Province   Province  `gorm:"foreignKey:ProvinceID;constraint:OnDelete:CASCADE;" json:"province,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Canton) TableName() string { return "cantons" }

// District representa un distrito
type District struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:text;not null" json:"name"`
	CantonID   uint      `gorm:"not null" json:"canton_id"`
	Canton     Canton    `gorm:"foreignKey:CantonID;constraint:OnDelete:CASCADE;" json:"canton,omitempty"`
	Latitude   float64   `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude  float64   `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	PostalCode string    `gorm:"type:text" json:"postal_code,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (District) TableName() string { return "districts" }
