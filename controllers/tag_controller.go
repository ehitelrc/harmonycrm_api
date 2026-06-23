package controllers

import (
	"harmony_api/dto"
	"harmony_api/models"
	"harmony_api/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	repo *repository.TagRepository
}

func NewTagController() *TagController {
	return &TagController{
		repo: repository.NewTagRepository(),
	}
}

func (c *TagController) Create(ctx *gin.Context) {
	var req dto.CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := &models.Tag{
		Name:         req.Name,
		Color:        req.Color,
		Icon:         req.Icon,
		DepartmentID: req.DepartmentID,
	}

	if err := c.repo.Create(tag); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating tag"})
		return
	}

	ctx.JSON(http.StatusCreated, tag)
}

func (c *TagController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	var req dto.UpdateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := c.repo.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}
	if req.Icon != "" {
		tag.Icon = req.Icon
	}
	if req.DepartmentID != nil {
		tag.DepartmentID = req.DepartmentID
	}

	if err := c.repo.Update(tag); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating tag"})
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

func (c *TagController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting tag"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

func (c *TagController) GetAll(ctx *gin.Context) {
	deptIDStr := ctx.Query("department_id")
	var tags []models.Tag
	var err error

	if deptIDStr != "" {
		deptID, parseErr := strconv.ParseUint(deptIDStr, 10, 32)
		if parseErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
			return
		}
		tags, err = c.repo.GetByDepartment(uint(deptID))
	} else {
		tags, err = c.repo.GetAll()
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving tags"})
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

// Relaciones con los casos
func (c *TagController) AssignToCase(ctx *gin.Context) {
	caseIdStr := ctx.Param("caseId")
	caseId, err := strconv.ParseUint(caseIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	var req dto.AssignTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.repo.AssignToCase(uint(caseId), req.TagID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error assigning tag"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tag assigned successfully"})
}

func (c *TagController) RemoveFromCase(ctx *gin.Context) {
	caseIdStr := ctx.Param("caseId")
	tagIdStr := ctx.Param("tagId")

	caseId, err1 := strconv.ParseUint(caseIdStr, 10, 32)
	tagId, err2 := strconv.ParseUint(tagIdStr, 10, 32)

	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameters"})
		return
	}

	if err := c.repo.RemoveFromCase(uint(caseId), uint(tagId)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing tag"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tag removed successfully"})
}

func (c *TagController) GetTagsByCase(ctx *gin.Context) {
	caseIdStr := ctx.Param("caseId")
	caseId, err := strconv.ParseUint(caseIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	tags, err := c.repo.GetTagsByCase(uint(caseId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving case tags"})
		return
	}

	ctx.JSON(http.StatusOK, tags)
}
