package models

import "time"

type CampaignReachSummary struct {
	SentCount      int64 `json:"sent_count"`
	DeliveredCount int64 `json:"delivered_count"`
	ReadCount      int64 `json:"read_count"`
	RepliedCount   int64 `json:"replied_count"`
	FailedCount    int64 `json:"failed_count"`
}

type CampaignReachRecipient struct {
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number"`
	FullName    string    `json:"full_name" gorm:"column:full_name"`
	CaseID      uint      `json:"case_id" gorm:"column:case_id"`
	Status      string    `json:"status" gorm:"column:status"`
	Replied     bool      `json:"replied" gorm:"column:replied"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
}

type CampaignReachResponse struct {
	Summary    CampaignReachSummary     `json:"summary"`
	Recipients []CampaignReachRecipient `json:"recipients"`
}
