package models

type Agent struct {
	UserID uint `json:"user_id" gorm:"primaryKey"`
}

func (Agent) TableName() string { return "agents" }
