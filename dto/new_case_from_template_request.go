package dto

type NewWhatsappCaseFromTemplateRequest struct {
	TemplateID           int    `json:"template_id" binding:"required"`
	ChannelIntegrationID uint   `json:"channel_integration_id" binding:"required"`
	ContactPhone         string `json:"contact_phone" binding:"required"`
	AgentID              int    `json:"agent_id" binding:"required"`
	ClientID             *uint  `json:"client_id" `
}
