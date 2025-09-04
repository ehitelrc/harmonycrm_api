package models

import "time"

type AgentDepartmentAssignment struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	AgentID      uint      `json:"agent_id"`      // FK -> agents.user_id
	DepartmentID uint      `json:"department_id"` // FK -> departments.id
	CreatedAt    time.Time `json:"created_at"`
}

func (AgentDepartmentAssignment) TableName() string { return "agent_department_assignments" }
