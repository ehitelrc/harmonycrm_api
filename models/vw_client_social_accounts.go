package models

import "time"

type VWClientSocialAccount struct {
	ClientID         uint      `json:"client_id"`
	ClientExternalID string    `json:"client_external_id"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	ClientCreatedAt  time.Time `json:"client_created_at"`
	ClientUpdatedAt  time.Time `json:"client_updated_at"`

	SocialAccountID  *uint      `json:"social_account_id"` // puede ser null si no tiene redes sociales
	ChannelID        *uint      `json:"channel_id"`
	SocialExternalID *string    `json:"social_external_id"`
	Username         *string    `json:"username"`
	ProfilePic       *string    `json:"profile_pic"`
	IsActive         *bool      `json:"is_active"`
	SocialCreatedAt  *time.Time `json:"social_created_at"`
	SocialUpdatedAt  *time.Time `json:"social_updated_at"`
}

func (VWClientSocialAccount) TableName() string {
	return "vw_client_social_accounts"
}
