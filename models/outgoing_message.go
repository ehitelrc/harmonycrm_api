package models

type AgentMessage struct {
	CaseID      uint   `json:"case_id"`
	SenderType  string `json:"sender_type"` // "agent" o "client"
	MessageType string `json:"message_type"`
	TextMessage string `json:"text_message"`
}
