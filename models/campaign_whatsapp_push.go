package models

import "time"

type CampaignWhatsappPush struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CampaignID  int64     `json:"campaign_id" gorm:"column:campaign_id;not null"`
	Description string    `json:"description" gorm:"column:description;size:50;not null"`
	TemplateID  int64     `json:"template_id" gorm:"column:template_id;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	ChangedBy   int       `json:"changed_by" gorm:"column:changed_by;index"`
}

// TableName especifica el nombre exacto de la tabla en la BD
func (CampaignWhatsappPush) TableName() string {
	return "campaign_whatsapp_push"
}
