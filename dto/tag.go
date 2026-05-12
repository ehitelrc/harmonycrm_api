package dto

type CreateTagRequest struct {
	Name      string `json:"name" binding:"required"`
	Color     string `json:"color" binding:"required"`
	Icon      string `json:"icon" binding:"required"`
}

type UpdateTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

type AssignTagRequest struct {
	TagID  uint `json:"tag_id" binding:"required"`
}
