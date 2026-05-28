package dto

import "time"

type BulkCloseCaseItem struct {
	CaseID                  uint       `json:"case_id"`
	ChannelLogoURL          string     `json:"channel_logo_url"`
	ChannelIntegrationName  string     `json:"channel_integration_name"`
	StartedAt               time.Time  `json:"started_at"`
	SenderId                string     `json:"sender_id"`
	ClientName              string     `json:"client_name"`
	MessagesSentCount       int        `json:"messages_sent_count"`
	MessagesReceivedCount   int        `json:"messages_received_count"`
	LastMessageDate         *time.Time `json:"last_message_date"`
	LastMessageType         string     `json:"last_message_type"` // 'sent' or 'received'
	PaymentFound            bool       `json:"payment_found"`
	PaymentRecordsCount     int        `json:"payment_records_count"`
}

type BulkCloseExecuteRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	CaseIDs  []uint `json:"case_ids" binding:"required"`
	Password string `json:"password" binding:"required"`
}
