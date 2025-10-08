package models

import "time"

// VwCaseItemsDetail representa la vista vw_case_items_detail en la base de datos.
// Esta vista combina información del caso, cliente, item y usuario que registró el producto.
type VwCaseItemsDetail struct {
	CaseItemID      int        `json:"case_item_id"`
	CaseID          int        `json:"case_id"`
	ItemID          int        `json:"item_id"`
	CompanyID       *int       `json:"company_id,omitempty"`
	DepartmentID    *int       `json:"department_id,omitempty"`
	CampaignID      *int       `json:"campaign_id,omitempty"`
	ClientID        *int       `json:"client_id,omitempty"`
	ClientName      *string    `json:"client_name,omitempty"`
	ClientEmail     *string    `json:"client_email,omitempty"`
	ClientPhone     *string    `json:"client_phone,omitempty"`
	ItemName        string     `json:"item_name"`
	ItemDescription *string    `json:"item_description,omitempty"`
	ItemType        string     `json:"item_type"`
	Price           float64    `json:"price"`
	Quantity        float64    `json:"quantity"`
	TotalAmount     float64    `json:"total_amount"`
	Acquired        bool       `json:"acquired"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedBy       *int       `json:"created_by,omitempty"`
	CreatedByName   *string    `json:"created_by_name,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	CaseStatus      *string    `json:"case_status,omitempty"`
	FunnelStage     *string    `json:"funnel_stage,omitempty"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	ClosedAt        *time.Time `json:"closed_at,omitempty"`
}

func (VwCaseItemsDetail) TableName() string { return "vw_case_items_detail" }
