package models

type AssignCaseToAgentRequest struct {
	CaseID    int `json:"case_id" binding:"required"`
	AgentID   int `json:"agent_id" binding:"required"`
	ChangedBy int `json:"changed_by"` // opcional: si usas auth, puede llenarse en el controlador
}
