package models

type AssignCaseToCampaignRequest struct {
	CaseID     int `json:"case_id" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	ChangedBy  int `json:"changed_by"` // opcional: si usas auth, puede llenarse en el controlador
}
