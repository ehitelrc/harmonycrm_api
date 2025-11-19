package models

type ChannelAgentClient struct {
	ID           int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	ChannelID    int64 `json:"channel_id"`
	AgentID      int64 `json:"agent_id"`
	ClientID     int64 `json:"client_id"`
	DepartmentID int64 `json:"department_id"`
}

func (ChannelAgentClient) TableName() string {
	return "channel_agent_client"
}
