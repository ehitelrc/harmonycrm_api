package models

type MoveCaseStagePayload struct {
	CaseID      int     `json:"case_id"`              // ID del caso
	FunnelID    int     `json:"funnel_id"`            // ID del funnel
	FromStageID *int    `json:"from_stage_id"`        // puede ser null
	ToStageID   int     `json:"to_stage_id"`          // nuevo stage
	Note        *string `json:"note,omitempty"`       // nota opcional
	ChangedBy   *int    `json:"changed_by,omitempty"` // usuario que hizo el cambio (opcional)
}
