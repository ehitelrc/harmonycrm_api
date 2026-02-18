package models

type AgentMessage struct {
	CaseID           uint   `json:"case_id"`
	SenderType       string `json:"sender_type"` // "agent" o "client"
	MessageType      string `json:"message_type"`
	TextMessage      string `json:"text_message"`
	Base64Content    string `json:"base64_content,omitempty"`
	MIMEType         string `json:"mime_type,omitempty"`
	HasError         bool   `json:"has_error"`
	MessageError     string `json:"message_error,omitempty"`
	FileName         string `json:"file_name,omitempty"`
	ChannelMessageID string `json:"channel_message_id,omitempty"`
}
