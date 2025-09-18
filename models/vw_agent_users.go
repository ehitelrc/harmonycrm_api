package models

type AgentUser struct {
	ID              int    `json:"id" gorm:"column:id"`
	Email           string `json:"email" gorm:"column:email"`
	FullName        string `json:"full_name" gorm:"column:full_name"`
	Phone           string `json:"phone" gorm:"column:phone"`
	IsActive        bool   `json:"is_active" gorm:"column:is_active"`
	ProfileImageUrl string `json:"profile_image_url" gorm:"column:profile_image_url"`
}

func (AgentUser) TableName() string {
	return "vw_agent_users" // nombre de la vista
}
