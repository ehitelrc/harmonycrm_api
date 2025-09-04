package models

import "time"

type ClientSocialAccount struct {
	ID         uint      `json:"id"`
	ClientID   uint      `json:"client_id"`
	ChannelID  uint      `json:"channel_id"`
	ExternalID string    `json:"external_id"`
	Username   string    `json:"username"`
	ProfilePic string    `json:"profile_pic"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *ClientSocialAccount) TableName() string {
	return "client_social_accounts"
}
