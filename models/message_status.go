package models

import "time"

type MessageStatus struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	ChannelMessageID string    `gorm:"not null" json:"channel_message_id"`
	MessageStatus    string    `gorm:"not null" json:"message_status"`
	Applied          bool      `gorm:"default:false" json:"applied"`
	CreatedAt        time.Time `json:"created_at" gorm:"default:now()"`
}

func (m *MessageStatus) TableName() string {
	return "message_status"
}
