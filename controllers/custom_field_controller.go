package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomFieldController struct {
	Repo *repository.CustomFieldRepository
}

func NewCustomFieldController(repo *repository.CustomFieldRepository) *CustomFieldController {
	return &CustomFieldController{Repo: repo}
}

// Save custom field value for an entity
func (ctl *CustomFieldController) SaveEntityCustomFieldValue(c *gin.Context) {

	var body []models.CustomFieldValuePayload

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON binding error", nil, err)
		return
	}

	println(body)

	error := ctl.Repo.SaveEntityCustomFieldValue(body)
	if error != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error saving custom field values", nil, error)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Custom field values saved successfully", nil, nil)

}
