package models

// CampaignWhatsappPushRequest representa el payload maestro-detalle
type CampaignWhatsappPushRequest struct {
	CampaignID  int64                           `json:"campaign_id" binding:"required"`
	Description string                          `json:"description" binding:"required"`
	TemplateID  int64                           `json:"template_id" binding:"required"`
	ChangedBy   int                             `json:"changed_by" binding:"required"`
	Leads       []CampaignWhatsappPushLeadInput `json:"leads" binding:"required"`
}

// CampaignWhatsappPushLeadInput representa un detalle del payload
type CampaignWhatsappPushLeadInput struct {
	PhoneNumber string  `json:"phone_number" binding:"required"`
	ClientID    *int64  `json:"client_id,omitempty"`
	CaseID      *int64  `json:"case_id,omitempty"`
	FullName    *string `json:"full_name,omitempty"`
	MessageSent bool    `json:"message_sent"`
}
