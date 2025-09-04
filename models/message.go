package models

import "time"

type Message struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CaseID        uint      `gorm:"not null" json:"case_id"`
	SenderType    string    `gorm:"type:enum('client','agent')" json:"sender_type"`
	MessageType   string    `gorm:"type:enum('text','image','file','audio')" json:"message_type"`
	TextContent   string    `gorm:"type:text" json:"text_content,omitempty"`
	FileURL       string    `gorm:"type:text" json:"file_url,omitempty"`
	MIMEType      string    `gorm:"type:text" json:"mime_type,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Base64Content string    `gorm:"type:text" json:"base64_content,omitempty"`
}

func (m *Message) TableName() string {
	return "messages"
}
