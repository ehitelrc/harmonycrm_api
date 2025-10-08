package models

import "time"

type CaseItem struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	CaseID    int       `gorm:"not null" json:"case_id"`
	ItemID    int       `gorm:"not null" json:"item_id"`
	Price     float64   `gorm:"type:numeric(14,2);not null;default:0" json:"price"`
	Quantity  float64   `gorm:"type:numeric(10,2);not null;default:1" json:"quantity"`
	Notes     *string   `gorm:"type:text" json:"notes,omitempty"`
	Acquired  bool      `gorm:"default:false" json:"acquired"`
	CreatedBy *int      `gorm:"index" json:"created_by,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relaciones
	Case        *Case `gorm:"foreignKey:CaseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"case,omitempty"`
	Item        *Item `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"item,omitempty"`
	CreatedUser *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:NO ACTION,OnDelete:SET NULL;" json:"created_user,omitempty"`
}

func (CaseItem) TableName() string { return "case_items" }
