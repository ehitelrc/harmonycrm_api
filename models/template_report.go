package models

import "time"

type BulkTemplateSend struct {
	PushID               int64      `json:"push_id"`
	Description          string     `json:"description"`
	AgentName            string     `json:"agent_name"`
	TemplateName         string     `json:"template_name"`
	CreatedAt            time.Time  `json:"created_at"`
	TotalRecipients      int64      `json:"total_recipients"`
	SuccessfulSends      int64      `json:"successful_sends"`
	FailedSends          int64      `json:"failed_sends"`
	DepartmentID         *int64     `json:"department_id"`
	DepartmentName       string     `json:"department_name"`
	ChannelIntegrationID int64      `json:"channel_integration_id"`
}

type IndividualTemplateSend struct {
	MessageID            int64      `json:"message_id"`
	CaseID               int64      `json:"case_id"`
	AgentName            string     `json:"agent_name"`
	TemplateName         string     `json:"template_name"`
	CreatedAt            time.Time  `json:"created_at"`
	ClientPhone          string     `json:"client_phone"`
	Status               string     `json:"status"`
	DepartmentID         *int64     `json:"department_id"`
	DepartmentName       string     `json:"department_name"`
	ChannelIntegrationID int64      `json:"channel_integration_id"`
}

type DepartmentSummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TemplateReportResponse struct {
	BulkSends       []BulkTemplateSend       `json:"bulk_sends"`
	IndividualSends []IndividualTemplateSend `json:"individual_sends"`
	Departments     []DepartmentSummary      `json:"departments"`
}
