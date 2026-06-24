package models

import (
	"database/sql/driver"
	"encoding/json"
)

// CompanyDashboard representa el resumen de métricas por compañía.
type CompanyDashboard struct {
	CompanyID         int64            `json:"company_id" gorm:"column:company_id"`
	CompanyName       string           `json:"company_name" gorm:"column:company_name"`
	TotalCases        int64            `json:"total_cases" gorm:"column:total_cases"`
	OpenCases         int64            `json:"open_cases" gorm:"column:open_cases"`
	ClosedCases       int64            `json:"closed_cases" gorm:"column:closed_cases"`
	ClosedToday       int64            `json:"closed_today" gorm:"column:closed_today"`
	OpenedToday       int64            `json:"opened_today" gorm:"column:opened_today"`
	CancelledCases    int64            `json:"cancelled_cases" gorm:"column:cancelled_cases"`
	UnansweredCases   int64            `json:"unanswered_cases" gorm:"column:unanswered_cases"`
	UnassignedAgents  int64            `json:"unassigned_agents" gorm:"column:unassigned_agents"`
	UnassignedClients int64            `json:"unassigned_clients" gorm:"column:unassigned_clients"`
	AvgCloseHours     float64          `json:"avg_close_hours" gorm:"column:avg_close_hours"`
	CasesByChannel    JSONCaseChannels `json:"cases_by_channel" gorm:"column:cases_by_channel;type:jsonb"`
	CasesByAgent      JSONCaseAgents   `json:"cases_by_agent" gorm:"column:cases_by_agent;type:jsonb"`
	OldestOpenCases   JSONOldestOpen   `json:"oldest_open_cases" gorm:"column:oldest_open_cases;type:jsonb"`
	DepartmentID      int64            `json:"department_id" gorm:"column:department_id;"`
}

// TableName especifica el nombre de la vista en PostgreSQL.
func (CompanyDashboard) TableName() string {
	return "vw_case_dashboard_by_company"
}

//
// --- Tipos JSON anidados ---
//

// CaseChannelStat representa las métricas por canal.
type CaseChannelStat struct {
	ChannelID   *int64  `json:"channel_id"`
	ChannelName *string `json:"channel_name"`
	OpenCases   int64   `json:"open_cases"`
	ClosedCases int64   `json:"closed_cases"`
}

// CaseAgentStat representa las métricas por agente.
type CaseAgentStat struct {
	AgentID       *int64   `json:"agent_id"`
	AgentName     *string  `json:"agent_name"`
	OpenCases     int64    `json:"open_cases"`
	ClosedCases   int64    `json:"closed_cases"`
	AvgCloseHours *float64 `json:"avg_close_hours"`
}

// OldestOpenCase representa los casos abiertos más antiguos.
type OldestOpenCase struct {
	CaseID                int64     `json:"case_id"`
	ClientName            *string   `json:"client_name"`
	ClientPhone           *string   `json:"client_phone"`
	CreatedAt             *SafeTime `json:"created_at"`
	LastMessageAt         *SafeTime `json:"last_message_at"`
	LastMessageSenderType *string   `json:"last_message_sender_type"`
}

//
// --- Tipos personalizados para JSONB ---
//

// JSONCaseChannels maneja el campo JSONB de casos por canal.
type JSONCaseChannels []CaseChannelStat

func (j *JSONCaseChannels) Scan(value interface{}) error {
	if value == nil {
		*j = []CaseChannelStat{}
		return nil
	}
	return json.Unmarshal(value.([]byte), j)
}

func (j JSONCaseChannels) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// JSONCaseAgents maneja el campo JSONB de casos por agente.
type JSONCaseAgents []CaseAgentStat

func (j *JSONCaseAgents) Scan(value interface{}) error {
	if value == nil {
		*j = []CaseAgentStat{}
		return nil
	}
	return json.Unmarshal(value.([]byte), j)
}

func (j JSONCaseAgents) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// JSONOldestOpen maneja el campo JSONB de los casos más antiguos.
type JSONOldestOpen []OldestOpenCase

func (j *JSONOldestOpen) Scan(value interface{}) error {
	if value == nil {
		*j = []OldestOpenCase{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok || len(bytes) == 0 {
		*j = []OldestOpenCase{}
		return nil
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONOldestOpen) Value() (driver.Value, error) {
	return json.Marshal(j)
}
