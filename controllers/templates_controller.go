package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"harmony_api/config"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"
	"strings"

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

	if template.ButtonsJSON != nil && *template.ButtonsJSON == "" {
		template.ButtonsJSON = nil
	}

	if err := c.repo.Create(&template); err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al crear la plantilla", nil, err)
		return
	}

	// 1. Obtener las integraciones activas de este canal (sin vincularlas automáticamente a la plantilla)
	var integrations []models.ChannelIntegration
	config.DB.Where("channel_id = ? AND is_active = ?", template.ChannelID, true).Find(&integrations)

	// 2. Intentar registrar automáticamente en Meta
	var wabaID string
	var accessToken string

	for _, integration := range integrations {
		if integration.MetaWabaID != nil && *integration.MetaWabaID != "" && integration.AccessToken != "" {
			wabaID = *integration.MetaWabaID
			accessToken = integration.AccessToken
			break
		}
	}

	// Fallback a WABA ID del canal o global, y primer token de acceso disponible
	if wabaID == "" && len(integrations) > 0 {
		var channel models.Channel
		if err := config.DB.Where("id = ?", template.ChannelID).First(&channel).Error; err == nil && channel.MetaWabaID != nil && *channel.MetaWabaID != "" {
			wabaID = *channel.MetaWabaID
		} else {
			wabaID, _ = repository.GetSettingTextValue("WAB_ID")
		}

		for _, integration := range integrations {
			if integration.AccessToken != "" {
				accessToken = integration.AccessToken
				break
			}
		}
	}

	if wabaID != "" && accessToken != "" && len(integrations) > 0 {
		if template.BodyContent != nil && *template.BodyContent != "" {
			fmt.Printf("🚀 Auto-registrando nueva plantilla '%s' en Meta...\n", template.TemplateName)
			metaID, status, err := RegisterTemplateInMeta(wabaID, accessToken, &template)
			if err != nil {
				fmt.Printf("⚠️ Error en auto-registro de plantilla '%s': %v\n", template.TemplateName, err)
			} else {
				statusMap := map[string]string{
					"APPROVED":         "approved",
					"PENDING_APPROVAL": "pending",
					"REJECTED":         "rejected",
					"PAUSED":           "paused",
					"DISABLED":         "disabled",
				}
				appStatus := "pending"
				if s, exists := statusMap[status]; exists {
					appStatus = s
				}

				config.DB.Model(&template).Updates(map[string]interface{}{
					"approval_status":  appStatus,
					"meta_template_id": metaID,
				})
				fmt.Printf("✅ Plantilla '%s' registrada en Meta automáticamente. ID: %s, Estado: %s\n", template.TemplateName, metaID, appStatus)
			}
		}
	} else if template.BodyContent != nil && *template.BodyContent != "" {
		// Simulación local si no tiene credenciales de Meta o si no hay integraciones vinculadas
		config.DB.Model(&template).Updates(map[string]interface{}{
			"approval_status":  "pending",
			"meta_template_id": "simulated_id_" + strconv.Itoa(int(template.ID)),
		})
	}

	// Recargar para obtener el estado e IDs actualizados
	updated, err := c.repo.GetByID(template.ID)
	if err == nil {
		template = *updated
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
		MetaTemplateID        *string `json:"meta_template_id"`
		MetaCategory          *string `json:"meta_category"`
		ApprovalStatus        *string `json:"approval_status"`
		RejectionReason       *string `json:"rejection_reason"`
		BodyContent           *string `json:"body_content"`
		HeaderContent         *string `json:"header_content"`
		FooterContent         *string `json:"footer_content"`
		ButtonsJSON           *string `json:"buttons_json"`
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
	if body.MetaTemplateID != nil {
		payload["meta_template_id"] = *body.MetaTemplateID
	}
	if body.MetaCategory != nil {
		payload["meta_category"] = *body.MetaCategory
	}
	if body.ApprovalStatus != nil {
		payload["approval_status"] = *body.ApprovalStatus
	}
	if body.RejectionReason != nil {
		payload["rejection_reason"] = *body.RejectionReason
	}
	if body.BodyContent != nil {
		payload["body_content"] = *body.BodyContent
	}
	if body.HeaderContent != nil {
		payload["header_content"] = *body.HeaderContent
	}
	if body.FooterContent != nil {
		payload["footer_content"] = *body.FooterContent
	}
	if body.ButtonsJSON != nil {
		if *body.ButtonsJSON == "" {
			payload["buttons_json"] = nil
		} else {
			payload["buttons_json"] = *body.ButtonsJSON
		}
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

	template, err := c.repo.GetByID(uint(id))
	if err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, "Plantilla no encontrada", nil, err)
		return
	}

	// Si está registrada en Meta, intentar eliminarla de Meta también
	if template.MetaTemplateID != nil && !strings.HasPrefix(*template.MetaTemplateID, "simulated_id_") {
		var wabaID string
		var accessToken string
		var integration models.ChannelIntegration
		if err := config.DB.Where("channel_id = ? AND is_active = ? AND access_token IS NOT NULL AND access_token != ''", template.ChannelID, true).First(&integration).Error; err == nil {
			accessToken = integration.AccessToken
			if integration.MetaWabaID != nil && *integration.MetaWabaID != "" {
				wabaID = *integration.MetaWabaID
			}
		}

		if wabaID == "" {
			var channel models.Channel
			if err := config.DB.Where("id = ?", template.ChannelID).First(&channel).Error; err == nil && channel.MetaWabaID != nil && *channel.MetaWabaID != "" {
				wabaID = *channel.MetaWabaID
			} else {
				wabaID, _ = repository.GetSettingTextValue("WAB_ID")
			}
		}

		if wabaID != "" && accessToken != "" {
			err = DeleteTemplateInMeta(wabaID, accessToken, template.TemplateName)
			if err != nil {
				fmt.Printf("⚠️ No se pudo eliminar plantilla '%s' de Meta: %v. Se eliminará localmente.\n", template.TemplateName, err)
			} else {
				fmt.Printf("✅ Plantilla '%s' eliminada de Meta con éxito.\n", template.TemplateName)
			}
		}
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

	// Obtener la integración para sacar el access token y waba id
	var integration models.ChannelIntegration
	if err := config.DB.Where("id = ?", integrationID).First(&integration).Error; err != nil {
		utils.Respond(ctx, http.StatusNotFound, false, "Integración no encontrada", nil, err)
		return
	}

	var wabaID string
	if integration.MetaWabaID != nil && *integration.MetaWabaID != "" {
		wabaID = *integration.MetaWabaID
	} else {
		wabaID, _ = repository.GetSettingTextValue("WAB_ID")
	}

	bodyText, err := repository.GetTemplateBodyFromMeta(templateName, wabaID, integration.AccessToken)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error consultando a Meta", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Plantilla obtenida de Meta", bodyText, nil)
}

func (c *TemplateController) RegisterMetaTemplate(ctx *gin.Context) {
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

	if template.BodyContent == nil || *template.BodyContent == "" {
		utils.Respond(ctx, http.StatusBadRequest, false, "El cuerpo de la plantilla (body_content) es requerido para registrar en Meta", nil, nil)
		return
	}

	// 1. Obtener integraciones activas del canal de la plantilla
	var integrations []models.ChannelIntegration
	if err := config.DB.Where("channel_id = ? AND is_active = ? AND access_token IS NOT NULL AND access_token != ''", template.ChannelID, true).Find(&integrations).Error; err != nil || len(integrations) == 0 {
		// Mock/Simulación de registro si no tiene integraciones configuradas en el canal (útil para pruebas locales)
		status := "pending"
		err = c.repo.Update(template.ID, map[string]interface{}{
			"approval_status":  status,
			"meta_template_id": "simulated_id_" + strconv.Itoa(int(template.ID)),
		})
		if err != nil {
			utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar estado local", nil, err)
			return
		}
		
		updated, _ := c.repo.GetByID(template.ID)
		utils.Respond(ctx, http.StatusOK, true, "Registro en Meta simulado con éxito (En Revisión)", updated, nil)
		return
	}

	// 2. Intentar registrar en Meta usando la primera integración activa del canal
	integration := integrations[0]

	var wabaID string
	if integration.MetaWabaID != nil && *integration.MetaWabaID != "" {
		wabaID = *integration.MetaWabaID
	}

	if wabaID == "" {
		var channel models.Channel
		if err := config.DB.Where("id = ?", template.ChannelID).First(&channel).Error; err == nil && channel.MetaWabaID != nil && *channel.MetaWabaID != "" {
			wabaID = *channel.MetaWabaID
		} else {
			wabaID, _ = repository.GetSettingTextValue("WAB_ID")
		}
	}

	if wabaID == "" || integration.AccessToken == "" {
		// Fallback a simulación si faltan credenciales reales
		status := "pending"
		err = c.repo.Update(template.ID, map[string]interface{}{
			"approval_status":  status,
			"meta_template_id": "simulated_id_" + strconv.Itoa(int(template.ID)),
		})
		if err != nil {
			utils.Respond(ctx, http.StatusInternalServerError, false, "Error al actualizar estado local", nil, err)
			return
		}
		updated, _ := c.repo.GetByID(template.ID)
		utils.Respond(ctx, http.StatusOK, true, "Registro en Meta simulado con éxito (Credenciales incompletas)", updated, nil)
		return
	}

	// 3. Registrar en Meta API real
	metaID, metaStatus, err := RegisterTemplateInMeta(wabaID, integration.AccessToken, template)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error registrando en Meta: "+err.Error(), nil, err)
		return
	}

	statusMap := map[string]string{
		"PENDING_APPROVAL": "pending",
		"APPROVED":         "approved",
		"REJECTED":         "rejected",
	}
	appStatus := "pending"
	if s, exists := statusMap[metaStatus]; exists {
		appStatus = s
	}

	err = c.repo.Update(template.ID, map[string]interface{}{
		"approval_status":  appStatus,
		"meta_template_id": metaID,
	})
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al guardar el ID de Meta", nil, err)
		return
	}

	updated, _ := c.repo.GetByID(template.ID)
	utils.Respond(ctx, http.StatusOK, true, "Registrada en Meta con éxito. Estado: "+metaStatus, updated, nil)
}

func (c *TemplateController) SyncMetaTemplate(ctx *gin.Context) {
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

	// 1. Simulación si es una plantilla simulada o si no tiene credenciales configuradas
	isSimulated := false
	if template.MetaTemplateID != nil && strings.HasPrefix(*template.MetaTemplateID, "simulated_id_") {
		isSimulated = true
	}

	// Buscar credenciales
	var wabaID string
	var accessToken string
	var integration models.ChannelIntegration
	if err := config.DB.Where("channel_id = ? AND is_active = ? AND access_token IS NOT NULL AND access_token != ''", template.ChannelID, true).First(&integration).Error; err == nil {
		accessToken = integration.AccessToken
		if integration.MetaWabaID != nil && *integration.MetaWabaID != "" {
			wabaID = *integration.MetaWabaID
		}
	}

	// Fallback
	if wabaID == "" {
		var channel models.Channel
		if err := config.DB.Where("id = ?", template.ChannelID).First(&channel).Error; err == nil && channel.MetaWabaID != nil && *channel.MetaWabaID != "" {
			wabaID = *channel.MetaWabaID
		} else {
			wabaID, _ = repository.GetSettingTextValue("WAB_ID")
		}
	}

	if isSimulated || wabaID == "" || accessToken == "" {
		// Mock local: marcar como aprobado y poblar datos si estaban en null
		status := "approved"
		if strings.Contains(template.TemplateName, "test_reject") {
			status = "rejected"
		}
		
		simID := "simulated_id_" + strconv.Itoa(int(template.ID))
		mockUpdates := map[string]interface{}{
			"approval_status":  status,
			"meta_template_id": simID,
			"meta_category":    "UTILITY",
		}
		
		if template.BodyContent == nil || *template.BodyContent == "" {
			bodyText := "Hola, este es un cuerpo de mensaje simulado de WhatsApp recuperado para la plantilla: " + template.TemplateName
			mockUpdates["body_content"] = bodyText
		}
		if template.HeaderContent == nil || *template.HeaderContent == "" {
			mockUpdates["header_content"] = "Encabezado Simulada"
		}
		if template.FooterContent == nil || *template.FooterContent == "" {
			mockUpdates["footer_content"] = "Pie de página simulado"
		}

		err = c.repo.Update(template.ID, mockUpdates)
		if err != nil {
			utils.Respond(ctx, http.StatusInternalServerError, false, "Error actualizando estado simulado", nil, err)
			return
		}
		updated, _ := c.repo.GetByID(template.ID)
		utils.Respond(ctx, http.StatusOK, true, "Sincronización simulada exitosa. Estado: "+status, updated, nil)
		return
	}

	// 2. Sincronización real con Meta
	updates, err := SyncAndFetchMetaTemplateDetails(wabaID, accessToken, template.TemplateName)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al sincronizar con Meta: "+err.Error(), nil, err)
		return
	}

	err = c.repo.Update(template.ID, updates)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error guardando la sincronización de Meta", nil, err)
		return
	}

	updated, _ := c.repo.GetByID(template.ID)
	utils.Respond(ctx, http.StatusOK, true, "Sincronizado con Meta con éxito. Estado: "+updates["approval_status"].(string), updated, nil)
}

type MetaComponent struct {
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
	Text   string `json:"text,omitempty"`
}

type MetaTemplateData struct {
	ID         string          `json:"id"`
	Status     string          `json:"status"`
	Category   string          `json:"category"`
	Components []MetaComponent `json:"components"`
}

type MetaTemplatesResponse struct {
	Data []MetaTemplateData `json:"data"`
}

func SyncAndFetchMetaTemplateDetails(wabaID string, accessToken string, templateName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/message_templates?name=%s", wabaID, templateName)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Meta API devolvió estado %s", res.Status)
	}

	var metaResp MetaTemplatesResponse
	if err := json.NewDecoder(res.Body).Decode(&metaResp); err != nil {
		return nil, err
	}

	if len(metaResp.Data) == 0 {
		return nil, fmt.Errorf("plantilla no encontrada en Meta")
	}

	target := metaResp.Data[0]
	
	statusMap := map[string]string{
		"APPROVED":         "approved",
		"PENDING_APPROVAL": "pending",
		"REJECTED":         "rejected",
		"PAUSED":           "paused",
		"DISABLED":         "disabled",
	}
	
	appStatus := "pending"
	if s, exists := statusMap[target.Status]; exists {
		appStatus = s
	}

	updates := map[string]interface{}{
		"approval_status":  appStatus,
		"meta_template_id": target.ID,
		"meta_category":    target.Category,
	}

	for _, comp := range target.Components {
		switch comp.Type {
		case "HEADER":
			if comp.Format == "TEXT" || comp.Format == "" {
				updates["header_content"] = comp.Text
			}
		case "BODY":
			updates["body_content"] = comp.Text
		case "FOOTER":
			updates["footer_content"] = comp.Text
		}
	}

	return updates, nil
}

func RegisterTemplateInMeta(wabaID string, accessToken string, template *models.MessageTemplate) (string, string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/message_templates", wabaID)

	type MetaTemplateComponent struct {
		Type    string        `json:"type"`
		Format  string        `json:"format,omitempty"`
		Text    string        `json:"text,omitempty"`
		Buttons []interface{} `json:"buttons,omitempty"`
	}

	type MetaTemplatePayload struct {
		Name       string                  `json:"name"`
		Category   string                  `json:"category"`
		Language   string                  `json:"language"`
		Components []MetaTemplateComponent `json:"components"`
	}

	var components []MetaTemplateComponent

	// 1. HEADER (if present)
	if template.HeaderContent != nil && *template.HeaderContent != "" {
		components = append(components, MetaTemplateComponent{
			Type:   "HEADER",
			Format: "TEXT",
			Text:   *template.HeaderContent,
		})
	}

	// 2. BODY (always required)
	bodyText := ""
	if template.BodyContent != nil {
		bodyText = *template.BodyContent
	}
	components = append(components, MetaTemplateComponent{
		Type: "BODY",
		Text: bodyText,
	})

	// 3. FOOTER (if present)
	if template.FooterContent != nil && *template.FooterContent != "" {
		components = append(components, MetaTemplateComponent{
			Type: "FOOTER",
			Text: *template.FooterContent,
		})
	}

	// 4. BUTTONS (if present)
	if template.ButtonsJSON != nil && *template.ButtonsJSON != "" {
		var buttons []interface{}
		if err := json.Unmarshal([]byte(*template.ButtonsJSON), &buttons); err == nil && len(buttons) > 0 {
			components = append(components, MetaTemplateComponent{
				Type:    "BUTTONS",
				Buttons: buttons,
			})
		} else if err != nil {
			fmt.Printf("⚠️ Error deserializando buttons_json para Meta: %v\n", err)
		}
	}

	metaCategory := "UTILITY"
	if template.MetaCategory != nil && *template.MetaCategory != "" {
		metaCategory = *template.MetaCategory
	}

	payload := MetaTemplatePayload{
		Name:       template.TemplateName,
		Category:   metaCategory,
		Language:   template.LanguageCode,
		Components: components,
	}

	jsonBytes, _ := json.Marshal(payload)
	fmt.Printf("📤 Meta Payload: %s\n", string(jsonBytes))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.NewDecoder(res.Body).Decode(&errResp)
		fmt.Printf("⚠️ Meta Error Response: %+v\n", errResp)
		errMsg := "Meta API returned status " + res.Status
		if errObj, ok := errResp["error"].(map[string]interface{}); ok {
			if msg, exists := errObj["message"].(string); exists {
				errMsg = msg
			}
		}
		return "", "", fmt.Errorf("%s", errMsg)
	}

	var metaResp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	json.NewDecoder(res.Body).Decode(&metaResp)

	return metaResp.ID, metaResp.Status, nil
}

func GetTemplateStatusFromMeta(wabaID string, accessToken string, templateName string) (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/message_templates?name=%s", wabaID, templateName)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("Meta API returned status %s", res.Status)
	}

	var metaResp struct {
		Data []struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&metaResp); err != nil {
		return "", err
	}

	if len(metaResp.Data) == 0 {
		return "", fmt.Errorf("template not found in Meta")
	}

	return metaResp.Data[0].Status, nil
}

func DeleteTemplateInMeta(wabaID string, accessToken string, templateName string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/message_templates?name=%s", wabaID, templateName)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		json.NewDecoder(res.Body).Decode(&errResp)
		if errResp.Error.Message != "" {
			return fmt.Errorf("Meta API: %s", errResp.Error.Message)
		}
		return fmt.Errorf("Meta API devolvió estado %s", res.Status)
	}

	return nil
}
