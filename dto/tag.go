package dto

type CreateTagRequest struct {
	Name         string `json:"name" binding:"required"`
	Color        string `json:"color" binding:"required"`
	Icon         string `json:"icon" binding:"required"`
	DepartmentID *uint  `json:"department_id"`
}

type UpdateTagRequest struct {
	Name         string `json:"name"`
	Color        string `json:"color"`
	Icon         string `json:"icon"`
	DepartmentID *uint  `json:"department_id"`
}

type AssignTagRequest struct {
	TagID  uint `json:"tag_id" binding:"required"`
}
