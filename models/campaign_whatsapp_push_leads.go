package models

// CampaignWhatsappPushLead representa los leads de un push de WhatsApp
type CampaignWhatsappPushLead struct {
	ID          int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	PushID      int64   `json:"push_id"`
	PhoneNumber string  `json:"phone_number"`
	ClientID    *int64  `json:"client_id,omitempty"` // puede ser NULL
	CaseID      *int64  `json:"case_id,omitempty"`   // puede ser NULL
	MessageSent bool    `json:"message_sent"`
	FullName    *string `json:"full_name,omitempty"` // opcional si agregás columna

}

func (CampaignWhatsappPushLead) TableName() string {
	return "campaign_whatsapp_push_leads"
}
