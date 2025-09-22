package models

type AssignCaseToDepartmentRequest struct {
	CaseID       int `json:"case_id" binding:"required"`
	DepartmentID int `json:"department_id" binding:"required"`
	ChangedBy    int `json:"changed_by"` // opcional: si usas auth, puede llenarse en el controlador
}
