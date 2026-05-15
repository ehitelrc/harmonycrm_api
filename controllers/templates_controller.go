package controllers

import (
	"harmony_api/config"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateController struct {
	repo *repository.TemplateRepository
}

func NewTemplateController() *TemplateController {
	return &TemplateController{
		repo: repository.NewTemplateRepository(),
	}
}

func (c *TemplateController) GetAllTemplates(ctx *gin.Context) {
	var channelID *uint
	if raw := ctx.Query("channel_id"); raw != "" {
		id, err := strconv.Atoi(raw)
		if err != nil {
			utils.Respond(ctx, http.StatusBadRequest, false, "channel_id inválido", nil, err)
			return
		}
		uid := uint(id)
		channelID = &uid
	}

	templates, err := c.repo.GetAll(channelID)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener las plantillas", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Lista de plantillas", templates, nil)
}

func (c *TemplateController) GetTemplateByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	template, err := c.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, "Plantilla no encontrada", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Plantilla encontrada", template, nil)
}

func (c *TemplateController) CreateTemplate(ctx *gin.Context) {
	var template models.MessageTemplate
	if err := ctx.ShouldBindJSON(&template); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	if err := c.repo.Create(&template); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al crear la plantilla", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusCreated, true, "Plantilla creada correctamente", template, nil)
}

func (c *TemplateController) UpdateTemplate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	var body struct {
		ChannelID             *uint   `json:"channel_id"`
		TemplateName          *string `json:"template_name"`
		LanguageCode          *string `json:"language_code"`
		Description           *string `json:"description"`
		Category              *string `json:"category"`
		IsActive              *bool   `json:"is_active"`
		IsConversationStarter *bool   `json:"is_conversation_starter"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	payload := map[string]interface{}{}
	if body.ChannelID != nil {
		payload["channel_id"] = *body.ChannelID
	}
	if body.TemplateName != nil {
		payload["template_name"] = *body.TemplateName
	}
	if body.LanguageCode != nil {
		payload["language_code"] = *body.LanguageCode
	}
	if body.Description != nil {
		payload["description"] = *body.Description
	}
	if body.Category != nil {
		payload["category"] = *body.Category
	}
	if body.IsActive != nil {
		payload["is_active"] = *body.IsActive
	}
	if body.IsConversationStarter != nil {
		payload["is_conversation_starter"] = *body.IsConversationStarter
	}

	if len(payload) == 0 {
		utils.Respond(ctx, http.StatusBadRequest, false, "No se enviaron campos para actualizar", nil, nil)
		return
	}

	if err := c.repo.Update(uint(id), payload); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar la plantilla", nil, err)
		return
	}

	// Return updated record
	updated, err := c.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusOK, true, "Plantilla actualizada correctamente", nil, nil)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Plantilla actualizada correctamente", updated, nil)
}

func (c *TemplateController) DeleteTemplate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al eliminar la plantilla", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Plantilla eliminada correctamente", nil, nil)
}

// ---- Integration-Template relations ----

func (c *TemplateController) GetTemplatesByIntegration(ctx *gin.Context) {
	integrationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de integración inválido", nil, err)
		return
	}

	templates, err := c.repo.GetByIntegrationID(uint(integrationID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener plantillas de la integración", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Plantillas de la integración", templates, nil)
}

func (c *TemplateController) GetIntegrationsForTemplate(ctx *gin.Context) {
	templateID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de plantilla inválido", nil, err)
		return
	}

	integrations, err := c.repo.GetIntegrationsForTemplate(uint(templateID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener integraciones de la plantilla", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Integraciones de la plantilla", integrations, nil)
}

func (c *TemplateController) CreateTemplateIntegration(ctx *gin.Context) {
	integrationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de integración inválido", nil, err)
		return
	}

	var body struct {
		TemplateID uint `json:"template_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	rel, err := c.repo.CreateIntegrationTemplate(uint(integrationID), body.TemplateID)
	if err != nil {
		utils.Respond(ctx, http.StatusConflict, false, err.Error(), nil, err)
		return
	}

	utils.Respond(ctx, http.StatusCreated, true, "Plantilla asignada a la integración correctamente", rel, nil)
}

func (c *TemplateController) DeleteTemplateIntegration(ctx *gin.Context) {
	integrationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de integración inválido", nil, err)
		return
	}

	templateIDStr := ctx.Query("template_id")
	if templateIDStr == "" {
		utils.Respond(ctx, http.StatusBadRequest, false, "template_id es requerido", nil, nil)
		return
	}
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "template_id inválido", nil, err)
		return
	}

	if err := c.repo.DeleteIntegrationTemplate(uint(integrationID), uint(templateID)); err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, err.Error(), nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Plantilla desasignada de la integración correctamente", nil, nil)
}

func (c *TemplateController) PreviewMetaTemplate(ctx *gin.Context) {
	templateName := ctx.Param("template_name")
	integrationIDParam := ctx.Param("integration_id")
	integrationID, err := strconv.Atoi(integrationIDParam)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de integración inválido", nil, err)
		return
	}

	// Obtener la integración para sacar el access token
	var integration models.ChannelIntegration
	if err := config.DB.Where("id = ?", integrationID).First(&integration).Error; err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, "Integración no encontrada", nil, err)
		return
	}

	bodyText, err := repository.GetTemplateBodyFromMeta(templateName, integration.AccessToken)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error consultando a Meta", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Plantilla obtenida de Meta", bodyText, nil)
}
