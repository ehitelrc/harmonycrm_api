package models

type AssignRequest struct {
	RoleID        uint `json:"role_id" binding:"required"`
	PermissionID  uint `json:"permission_id" binding:"required"`
	AssignRequest bool `json:"assign_request"` // para evitar warning de campo no usado
}
