package models

import "time"

type CaseNote struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CaseID    uint      `gorm:"not null;index" json:"case_id"`
	AuthorID  uint      `gorm:"not null;index" json:"author_id"`
	Note      string    `gorm:"type:text;not null" json:"note"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relaciones opcionales
	Case   *Case `gorm:"foreignKey:CaseID;references:ID" json:"case,omitempty"`
	Author *User `gorm:"foreignKey:AuthorID;references:ID" json:"author,omitempty"`
}
