package dto

import "harmony_api/models"

type ClosedCaseResponse struct {
	Case            models.Case `json:"case" gorm:"embedded"`
	AgentName       string      `json:"agent_name"`
	IntegrationName string      `json:"integration_name"`
	ChannelType     string      `json:"channel_type"`
	Icon            string      `json:"icon"`
}
