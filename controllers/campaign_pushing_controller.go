package controllers

import (
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"harmony_api/ws"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CampaignPushingController struct {
	repo *repository.CampaignPushingRepository
	hub  *ws.Hub
}

func NewCampaignPushingController(hub *ws.Hub) *CampaignPushingController {
	return &CampaignPushingController{
		repo: repository.NewCampaignPushingRepository(),
		hub:  hub,
	}
}

// RegisterWhatsappCampaignPush registra un nuevo push de campaña con leads
func (ctrl *CampaignPushingController) RegisterWhatsappCampaignPush(c *gin.Context) {
	var requestBody models.CampaignWhatsappPushRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Invalid request body", nil, err)
		return
	}

	// Guardar en DB
	pushID, err := ctrl.repo.CreateWhatsappPush(&requestBody, ctrl.hub)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error saving whatsapp campaign push", nil, err)
		return
	}

	// Respuesta
	utils.Respond(c, http.StatusOK, true, "Whatsapp campaign push created successfully", map[string]interface{}{
		"push_id": pushID,
	}, nil)
}
