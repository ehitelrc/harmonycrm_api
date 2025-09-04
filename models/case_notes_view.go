package models

import "time"

// Vista: public.case_notes_view
type CaseNoteView struct {
	ID          int        `json:"id"            gorm:"column:id"`
	CaseID      *int       `json:"case_id"       gorm:"column:case_id"`
	AuthorID    *int       `json:"author_id"     gorm:"column:author_id"`
	AuthorName  *string    `json:"author_name"   gorm:"column:author_name"`
	AuthorEmail *string    `json:"author_email"  gorm:"column:author_email"`
	Note        *string    `json:"note"          gorm:"column:note"`
	CreatedAt   *time.Time `json:"created_at"    gorm:"column:created_at"`
}

func (CaseNoteView) TableName() string { return "case_notes_view" }
