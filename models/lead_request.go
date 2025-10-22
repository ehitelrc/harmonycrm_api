// models/lead_request.go
package models

type LeadRequest struct {
	ClientID             uint            `json:"client_id"`
	CompanyID            uint            `json:"company_id"`
	CampaignID           uint            `json:"campaign_id"`
	FunnelStageID        uint            `json:"funnel_stage_id"`
	ChannelID            uint            `json:"channel_id"`
	ChannelIntegrationID uint            `json:"channel_integration_id"`
	AgentID              uint            `json:"agent_id"`
	Items                []ItemSelection `json:"items"`
}

type ItemSelection struct {
	ItemID    uint    `json:"item_id"`
	ItemName  string  `json:"item_name"`
	Quantity  int     `json:"quantity"`
	ItemPrice float64 `json:"item_price"`
	Notes     string  `json:"notes,omitempty"`
}
