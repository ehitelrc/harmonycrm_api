package models

type AgentDepartmentInformation struct {
	AgentID         uint   `json:"agent_id" gorm:"column:agent_id"`
	AgentName       string `json:"agent_name" gorm:"column:agent_name"`
	CompanyID       uint   `json:"company_id" gorm:"column:company_id"`
	DepartmentName  string `json:"department_name" gorm:"column:department_name"`
	DepartmentID    uint   `json:"department_id" gorm:"column:department_id"`
	ProfileImageURL string `json:"profile_image_url" gorm:"column:profile_image_url"`
}

func (AgentDepartmentInformation) TableName() string {
	return "vw_agent_department_information"
}
