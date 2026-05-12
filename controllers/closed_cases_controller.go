package controllers

import (
	"net/http"
	"strconv"

	"harmony_api/repository"
	"harmony_api/utils"

	"github.com/gin-gonic/gin"
)

type ClosedCasesController struct {
	repo *repository.MessageRepository
}

func NewClosedCasesController(repo *repository.MessageRepository) *ClosedCasesController {
	return &ClosedCasesController{repo: repo}
}

// Get closed cases for a given sender ID
func (ctrl *ClosedCasesController) GetClosedCases(c *gin.Context) {
	senderID := c.Param("sender_id")
	channelIntegrationID := c.Query("channel_integration_id")

	if senderID == "" {
		utils.Respond(c, http.StatusBadRequest, false, "Debe proveer un sender_id", nil, nil)
		return
	}

	cases, err := ctrl.repo.GetClosedCasesBySenderID(senderID, channelIntegrationID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error obteniendo los casos cerrados", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Casos cerrados obtenidos correctamente", cases, nil)
}

// Get messages for a given case ID specifically for closed cases viewer
func (ctrl *ClosedCasesController) GetCaseMessages(c *gin.Context) {
	caseIDStr := c.Param("case_id")

	if caseIDStr == "" {
		utils.Respond(c, http.StatusBadRequest, false, "Debe proveer un case_id", nil, nil)
		return
	}

	caseID, err := strconv.ParseUint(caseIDStr, 10, 32)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Formato de case_id inválido", nil, err)
		return
	}

	messages, err := ctrl.repo.GetClosedCaseMessages(uint(caseID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error obteniendo los mensajes del caso", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Mensajes del caso obtenidos correctamente", messages, nil)
}
