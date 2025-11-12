package models

import (
	"time"
)

type ChannelIntegration struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID     uint      `json:"company_id" gorm:"not null"`
	ChannelID     uint      `json:"channel_id" gorm:"not null"`
	WebhookURL    string    `json:"webhook_url" gorm:"type:text;not null"`
	AccessToken   string    `json:"access_token,omitempty" gorm:"type:text"`
	AppIdentifier string    `json:"app_identifier,omitempty" gorm:"type:text"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones opcionales (útiles para consultas extendidas o joins)
	Company         *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID;references:ID"`
	Channel         *Channel `json:"channel,omitempty" gorm:"foreignKey:ChannelID;references:ID"`
	IsNonCommercial bool     `gorm:"default:false" json:"is_non_commercial"`

	IntegrationName string `json:"integration_name,omitempty"` // Campo calculado, no se mapea a la base de datos

	DepartmentID *uint `json:"department_id,omitempty" gorm:"index"` // Índice para búsquedas rápidas por departamento
}

// Nombre explícito de la tabla (si tu esquema no usa pluralización)
func (ChannelIntegration) TableName() string {
	return "channel_integrations"
}
