package models

type VwAgentDepartmentAssignment struct {
	CompanyID          int    `json:"company_id"`
	DepartmentID       int    `json:"department_id"`
	DepartmentName     string `json:"department_name"`
	UserID             int    `json:"user_id"`
	DepartmentAssigned bool   `json:"department_assigned"`
}

func (VwAgentDepartmentAssignment) TableName() string {
	return "vw_agent_department_assignments"
}
