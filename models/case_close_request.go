package models

type CaseCloseRequest struct {
	CaseID   int    `json:"case_id"`        // ID del caso a cerrar
	FunnelID int    `json:"funnel_id"`      // ID del funnel asociado
	ClosedBy int    `json:"closed_by"`      // ID del usuario que cierra el caso
	Note     string `json:"note,omitempty"` // Razón para cerrar el caso (opcional)
}
