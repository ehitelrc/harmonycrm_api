// file: models/funnel.go
package models

// Funnel representa un embudo (pipeline) genérico.
// Tabla: funnels
type Funnel struct {
	ID          uint    `json:"id"           gorm:"primaryKey"`
	Name        string  `json:"name"         gorm:"size:120;not null"`
	Description *string `json:"description"  gorm:"type:text"`
	IsActive    bool    `json:"is_active"    gorm:"not null;default:true"`
}

// FunnelStage representa una etapa dentro de un embudo.
// Tabla: funnel_stages
type FunnelStage struct {
	ID       uint   `json:"id"        gorm:"primaryKey"`
	FunnelID uint   `json:"funnel_id" gorm:"not null;index"`
	Name     string `json:"name"      gorm:"size:120;not null"`

	// Orden dentro del funnel (0,1,2,...) — único por funnel
	Position int  `json:"position" gorm:"not null;index"`
	IsWon    bool `json:"is_won"    gorm:"not null;default:false"`
	IsLost   bool `json:"is_lost"   gorm:"not null;default:false"`
}

// Opcional: si prefieres declarar los nombres de tabla explícitos
func (Funnel) TableName() string      { return "funnels" }
func (FunnelStage) TableName() string { return "funnel_stages" }
