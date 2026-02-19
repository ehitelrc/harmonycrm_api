package repository

import (
	"encoding/json"
	"fmt"
	"harmony_api/config"
	"harmony_api/dto"
	"harmony_api/models"
	"harmony_api/utils"
	"harmony_api/ws"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

// Sen

type CampaignPushingRepository struct{}

func NewCampaignPushingRepository() *CampaignPushingRepository {
	return &CampaignPushingRepository{}
}

func (r *CampaignPushingRepository) CreateWhatsappPush(data *models.CampaignWhatsappPushRequest, hub *ws.Hub) (int64, error) {
	db := config.DB
	var pushID int64

	err := db.Transaction(func(tx *gorm.DB) error {

		// ===========================================================
		// 1. Crear encabezado solo si CampaignID != nil
		// ===========================================================
		if data.CampaignID != nil {
			header := models.CampaignWhatsappPush{
				CampaignID:  data.CampaignID,
				Description: data.Description,
				TemplateID:  data.TemplateID,
				ChangedBy:   data.ChangedBy,
			}

			if err := tx.Create(&header).Error; err != nil {
				return err
			}

			pushID = header.ID
		}

		// Buscar info del template
		var template models.CompanyChannelTemplateView
		if err := tx.Where("template_id = ?", data.TemplateID).First(&template).Error; err != nil {
			return err
		}

		// Listas donde meteremos leads + recipients para el envío masivo
		var leadsToInsert []models.CampaignWhatsappPushLead
		var recipients []models.TemplateRecipient

		// ===========================================================
		// 2. Procesar cada lead
		// ===========================================================
		for _, l := range data.Leads {

			// Primero buscar si existe un case abierto para este número
			var existingCase models.Case
			err := tx.
				Where("channel_integration_id = ? AND sender_id = ? AND status = ?",
					template.ChannelIntegrationID, l.PhoneNumber, "open").
				First(&existingCase).Error

			var caseID int64

			if err == nil {
				// Caso ya existe
				caseID = int64(existingCase.ID)

			} else if err == gorm.ErrRecordNotFound {

				// =======================================================
				// 2A. Crear un nuevo caso porque NO existe uno abierto
				// =======================================================

				// Para crear el caso ocupamos datos del integration
				var integration models.ViewChannelIntegration
				if err := tx.Where("channel_integration_id = ?", template.ChannelIntegrationID).First(&integration).Error; err != nil {
					return err
				}

				channelIDStr := strconv.FormatUint(uint64(integration.ChannelID), 10)

				newCase := models.Case{
					SenderId:             l.PhoneNumber,
					ChannelID:            channelIDStr,
					CompanyID:            integration.CompanyID,
					ChannelIntegrationID: integration.ChannelIntegrationID,
					IsNonCommercial:      integration.IsNonCommercial,
					DepartmentID:         *integration.DepartmentID,
					ClientID:             Int64PtrToUintPtr(l.ClientID),
					AgentID:              uint(data.ChangedBy),
					Status:               "open",
				}

				if err := tx.Create(&newCase).Error; err != nil {
					return err
				}

				bodyText, err := GetTemplateBodyFromMeta(*template.TemplateName, *template.AccessToken)
				if err != nil {
					bodyText = "Apertura mediante template"
				}

				// Crear mensaje inicial
				openMsg := models.Message{
					CaseID:      newCase.ID,
					SenderType:  "agent",
					MessageType: "text",
					TextContent: bodyText,
					MessageRead: true,
				}

				if err := tx.Create(&openMsg).Error; err != nil {
					return err
				}

				caseID = int64(newCase.ID)

			} else {
				// Error inesperado al buscar case
				return err
			}

			// Crear lead solo si CampaignID != nil (Campaña activa)
			if data.CampaignID != nil {
				lead := models.CampaignWhatsappPushLead{
					PushID:      pushID,
					PhoneNumber: l.PhoneNumber,
					ClientID:    l.ClientID,
					CaseID:      &caseID,
					FullName:    l.FullName,
					MessageSent: false,
				}
				leadsToInsert = append(leadsToInsert, lead)
			}

			// Agregar a recipients para envío del template
			recipients = append(recipients, models.TemplateRecipient{
				Number: l.PhoneNumber,
				CaseID: &caseID,
			})
		}

		// ===========================================================
		// 3. Insertar leads si corresponde
		// ===========================================================
		if len(leadsToInsert) > 0 {
			if err := tx.Create(&leadsToInsert).Error; err != nil {
				return err
			}
		}

		// ===========================================================
		// 4. Enviar mensaje de plantilla
		// ===========================================================
		utils.SendTemplateToMany(
			template.TemplateUrlWebhook,
			*template.AppIdentifier,
			*template.AccessToken,
			*template.TemplateName,
			*template.Language,
			recipients,
			hub,
		)

		return nil
	})

	if err != nil {
		return 0, err
	}

	return pushID, nil
}

// SendWhatsappTemplateMessage envía un mensaje de plantilla de WhatsApp para un caso específico
func (r *CampaignPushingRepository) SendWhatsappTemplateMessage(templateID int, caseID int) (*models.Message, error) {
	// Buscar template
	var template models.ChannelWhatsAppTemplate

	if err := config.DB.Where("id = ?", templateID).First(&template).Error; err != nil {
		return nil, err
	}

	// Buscar caso
	var caseChannelIntegration models.VWCaseChannelIntegration

	if err := config.DB.Where("case_id = ?", caseID).First(&caseChannelIntegration).Error; err != nil {
		return nil, err
	}

	// Enviar mensaje de plantilla
	caseIDParam := int64(caseChannelIntegration.CaseID)

	recipients := []models.TemplateRecipient{
		{
			Number: *caseChannelIntegration.SenderID,
			CaseID: &caseIDParam,
		},
	}

	// create message record

	bodyText, err := GetTemplateBodyFromMeta(*&template.TemplateName, *caseChannelIntegration.AccessToken)
	if err != nil {
		bodyText = "Apertura mediante template"
	}

	newMessage := models.Message{
		CaseID:      uint(caseIDParam),
		SenderType:  "agent",
		MessageType: "text",
		TextContent: bodyText,
		MessageRead: true,
	}

	if err := config.DB.Create(&newMessage).Error; err != nil {
		return nil, err
	}

	utils.SendTemplateToMany(template.TemplateUrlWebhook, *caseChannelIntegration.AppIdentifier, *caseChannelIntegration.AccessToken, *&template.TemplateName, *&template.Language, recipients, nil)

	return &newMessage, nil
}

func (r *CampaignPushingRepository) NewCaseFromTemplate(request dto.NewWhatsappCaseFromTemplateRequest, hub *ws.Hub) (int64, error) {
	db := config.DB

	var caseID int64

	err := db.Transaction(func(tx *gorm.DB) error {

		// Search channel ID by template ID
		var template models.CompanyChannelTemplateView

		if err := tx.Debug().Where("template_id = ?", request.TemplateID).First(&template).Error; err != nil {
			return err
		}

		var integration models.ViewChannelIntegration

		if err := tx.Where("channel_integration_id = ?", request.ChannelIntegrationID).First(&integration).Error; err != nil {
			return err
		}

		// Get case by channel_integration_id & sender_id and status open
		var newCase models.Case

		if err := tx.Debug().Where("channel_integration_id = ? AND sender_id = ? AND status = ?", request.ChannelIntegrationID, request.ContactPhone, "open").First(&newCase).Error; err != nil {
			if err == gorm.ErrRecordNotFound {

				// Create new case

				number := strconv.FormatUint(uint64(integration.ChannelID), 10)

				newCase := models.Case{
					SenderId:             request.ContactPhone,
					ChannelID:            number,
					CompanyID:            integration.CompanyID,
					ChannelIntegrationID: integration.ChannelIntegrationID,
					IsNonCommercial:      integration.IsNonCommercial,
					DepartmentID:         *integration.DepartmentID,
					ClientID:             request.ClientID,
					AgentID:              uint(request.AgentID),
					Status:               "open",
				}

				if err := tx.Create(&newCase).Error; err != nil {
					return err
				}

				bodyText, err := GetTemplateBodyFromMeta(*template.TemplateName, *template.AccessToken)
				if err != nil {
					bodyText = "Apertura mediante template"
				}

				newMessage := models.Message{
					CaseID:      newCase.ID,
					SenderType:  "agent",
					MessageType: "text",
					TextContent: bodyText,
					MessageRead: true,
				}

				if err := tx.Create(&newMessage).Error; err != nil {
					return err
				}

				caseID = int64(newCase.ID)

				// recipients := []models.TemplateRecipient{
				// 	{
				// 		Number: request.ContactPhone,
				// 		CaseID: &caseID,
				// 	},
				// }

				// utils.SendTemplateToMany(template.TemplateUrlWebhook, *template.AppIdentifier, *template.AccessToken, *template.TemplateName, *template.Language, recipients, hub)

				// ------------------------------------------------------------------
				// NUEVO: Enviar template individual y capturar ID (Wamid)
				// ------------------------------------------------------------------
				wamid, err := utils.SendTemplateMessageWithID(
					template.TemplateUrlWebhook,
					*template.AppIdentifier,
					*template.AccessToken,
					*template.TemplateName,
					*template.Language,
					request.ContactPhone,
				)

				if err != nil {
					return fmt.Errorf("error enviando template: %w", err)
				}

				// Actualizar el mensaje con el ID de Meta
				if err := tx.Model(&models.Message{}).Where("id = ?", newMessage.ID).Update("channel_message_id", wamid).Error; err != nil {
					fmt.Printf("⚠️ Error actualizando channel_message_id para mensaje %d: %v\n", newMessage.ID, err)
					// No retornamos error fatal para no hacer rollback de todo el caso
				}

			} else {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return caseID, nil
}

func Int64PtrToUintPtr(v *int64) *uint {
	if v == nil {
		return nil
	}
	u := uint(*v)
	return &u
}

func GetSettingTextValue(code string) (string, error) {
	var setting models.Setting
	if err := config.DB.Where("value_code = ? AND is_active = ?", code, true).First(&setting).Error; err != nil {
		return "", err
	}

	if setting.TextValue == nil {
		return "", fmt.Errorf("setting %s has NULL text_value", code)
	}

	return *setting.TextValue, nil
}

func GetTemplateBodyFromMeta(templateName string, accessToken string) (string, error) {
	wabaID, err := GetSettingTextValue("WAB_ID")
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://graph.facebook.com/v24.0/%s/message_templates?name=%s",
		wabaID,
		templateName,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var metaResp struct {
		Data []struct {
			Components []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"components"`
		} `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&metaResp); err != nil {
		return "", err
	}

	if len(metaResp.Data) == 0 {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	// Buscar el BODY
	for _, c := range metaResp.Data[0].Components {
		if c.Type == "BODY" {
			return c.Text, nil
		}
	}

	return "", fmt.Errorf("template BODY not found")
}
