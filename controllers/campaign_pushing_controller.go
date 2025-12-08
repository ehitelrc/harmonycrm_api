package controllers

import (
	"encoding/json"
	"harmony_api/dto"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"harmony_api/ws"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CampaignPushingController struct {
	campRepo *repository.ChannelRepository
	repo     *repository.CampaignPushingRepository
	hub      *ws.Hub
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

// SendWhatsappTemplateMessage envía un mensaje de plantilla de WhatsApp para un caso específico
func (ctrl *CampaignPushingController) SendWhatsappTemplateMessage(c *gin.Context) {

	templateID, err := strconv.Atoi(c.Param("template_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Invalid template ID", nil, err)
		return
	}

	caseID, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Invalid case ID", nil, err)
		return
	}

	err = ctrl.repo.SendWhatsappTemplateMessage(templateID, caseID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error sending WhatsApp template message", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "WhatsApp template message sent successfully", nil, nil)
}

func (ctrl *CampaignPushingController) CreateNewWhatsappCaseFromTemplate(c *gin.Context) {
	var requestBody dto.NewWhatsappCaseFromTemplateRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Invalid request body", nil, err)
		return
	}

	caseID, err := ctrl.repo.NewCaseFromTemplate(requestBody, ctrl.hub)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error creating new WhatsApp case from template", nil, err)
		return
	}
	// Broadcast WS (si tenemos case_id)
	if caseID != 0 && ctrl.hub != nil {
		payload, _ := json.Marshal(WSMessage{
			Type:   "new_message",
			CaseID: uint(caseID),
			Data:   "", // o arma un DTO si prefieres
		})
		channel := "case:" + strconv.Itoa(int(caseID))
		ctrl.hub.BroadcastJSON(channel, payload)

		ctrl.hub.BroadcastJSON("agent:"+strconv.Itoa(int(*&requestBody.AgentID)), payload)

	}

	// utils.Respond(c, http.StatusOK, true, "New WhatsApp case created successfully", map[string]interface{}{
	// 	"case_id": caseID,
	// }, nil)

	utils.Respond(c, http.StatusOK, true, "New WhatsApp case created successfully", map[string]interface{}{
		"case_id": caseID,
	}, nil)
}
