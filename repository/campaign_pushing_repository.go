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
	"strings"
	"sync"

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

	// Para crear el caso ocupamos datos del integration
	var integration models.ViewChannelIntegration

	if err := db.Where("channel_integration_id = ?", data.ChannelIntegrationID).First(&integration).Error; err != nil {
		return 0, err
	}

	// Buscar info del template
	var template models.ChannelTemplateIntegration
	if err := db.Where("template_id = ?", data.TemplateID).First(&template).Error; err != nil {
		return 0, err
	}

	// Obtener el cuerpo de la plantilla desde Meta una única vez fuera del bucle
	var wabaID string
	if integration.MetaWabaID != nil && *integration.MetaWabaID != "" {
		wabaID = *integration.MetaWabaID
	} else {
		wabaID, _ = GetSettingTextValue("WAB_ID")
	}

	templateBodyText, err := GetTemplateBodyFromMeta(template.TemplateName, wabaID, integration.AccessToken)
	if err != nil {
		fmt.Printf("⚠️ Error obteniendo cuerpo de plantilla desde Meta: %v. Usando texto por defecto.\n", err)
		templateBodyText = "Apertura mediante template"
	}

	// Listas donde meteremos leads + recipients para el envío masivo
	var leadsToInsert []models.CampaignWhatsappPushLead
	var recipients []models.TemplateRecipient
	messageIDs := make(map[string]uint)

	err = db.Transaction(func(tx *gorm.DB) error {

		// ===========================================================
		// 1. Crear encabezado solo si CampaignID != nil
		// ===========================================================
		if data.CampaignID != nil {
			header := models.CampaignWhatsappPush{
				CampaignID:           data.CampaignID,
				Description:          data.Description,
				TemplateID:           data.TemplateID,
				ChangedBy:            data.ChangedBy,
				ChannelIntegrationID: data.ChannelIntegrationID,
			}

			if err := tx.Create(&header).Error; err != nil {
				return err
			}

			pushID = header.ID
		}

		// ===========================================================
		// 2. Procesar cada lead
		// ===========================================================
		for _, l := range data.Leads {

			l.PhoneNumber = strings.TrimPrefix(l.PhoneNumber, "+")

			// Primero buscar si existe un case abierto para este número
			var existingCase models.Case
			err := tx.
				Where("channel_integration_id = ? AND sender_id = ? AND status = ?",
					template.IntegrationID, l.PhoneNumber, "open").
				First(&existingCase).Error

			var caseID int64

			switch err {
			case nil:
				// Caso ya existe
				caseID = int64(existingCase.ID)

			case gorm.ErrRecordNotFound:

				// =======================================================
				// 2A. Crear un nuevo caso porque NO existe uno abierto
				// =======================================================

				channelIDStr := strconv.FormatUint(uint64(integration.ChannelID), 10)

				finalAgentID := uint(data.ChangedBy)
				agentAssigned := false

				// Intento 1: Asignación directa en ChannelAgentClient
				if l.ClientID != nil {
					var channelAgentClient models.ChannelAgentClient
					errAgent := tx.Where("department_id = ? AND client_id = ?", *integration.DepartmentID, *l.ClientID).First(&channelAgentClient).Error
					if errAgent == nil {
						finalAgentID = uint(channelAgentClient.AgentID)
						agentAssigned = true
					}
				}

				// Intento 2: Historial del cliente en el canal (último caso atendido)
				if !agentAssigned {
					var lastCase models.Case
					errLastCase := tx.Where("channel_integration_id = ? AND sender_id = ?", integration.ChannelIntegrationID, l.PhoneNumber).Order("id desc").First(&lastCase).Error
					if errLastCase == nil {
						// Verificar si el agente asignado a ese caso previo sigue activo
						var agent models.User
						errUser := tx.Where("id = ?", lastCase.AgentID).First(&agent).Error
						if errUser == nil && agent.IsActive {
							finalAgentID = lastCase.AgentID
						}
					}
				}

				newCase := models.Case{
					SenderId:             l.PhoneNumber,
					ChannelID:            channelIDStr,
					CompanyID:            integration.CompanyID,
					ChannelIntegrationID: integration.ChannelIntegrationID,
					IsNonCommercial:      integration.IsNonCommercial,
					DepartmentID:         *integration.DepartmentID,
					ClientID:             Int64PtrToUintPtr(l.ClientID),
					AgentID:              finalAgentID,
					Status:               "open",
				}

				if err := tx.Create(&newCase).Error; err != nil {
					return err
				}

				// Crear mensaje inicial
				openMsg := models.Message{
					CaseID:      newCase.ID,
					SenderType:  "agent",
					MessageType: "text",
					TextContent: templateBodyText,
					MessageRead: true,
					TemplateID:  &template.TemplateID,
				}

				if err := tx.Create(&openMsg).Error; err != nil {
					return err
				}

				messageIDs[l.PhoneNumber] = openMsg.ID

				caseID = int64(newCase.ID)

			default:
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

		return nil
	})

	if err != nil {
		return 0, err
	}

	// ===========================================================
	// 4. Enviar mensajes de plantilla en segundo plano (asíncrono)
	//    con concurrencia controlada para evitar bloquear la API
	// ===========================================================
	go func(recipients []models.TemplateRecipient, messageIDs map[string]uint, integration models.ViewChannelIntegration, template models.ChannelTemplateIntegration) {
		dbConn := config.DB

		// Canal para distribuir los trabajos de envío
		jobs := make(chan models.TemplateRecipient, len(recipients))
		for _, r := range recipients {
			jobs <- r
		}
		close(jobs)

		// Limitar la concurrencia a un máximo de 15 trabajadores paralelos
		numWorkers := 15
		if len(recipients) < numWorkers {
			numWorkers = len(recipients)
		}

		var wg sync.WaitGroup
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for r := range jobs {
					if integration.ChannelCode != nil && *integration.ChannelCode == "whatsapp" {
						wamid, err := utils.SendTemplateMessageWithID(
							"https://graph.facebook.com/v24.0",
							integration.AppIdentifier,
							integration.AccessToken,
							template.TemplateName,
							template.LanguageCode,
							r.Number,
						)

						if err != nil {
							fmt.Printf("⚠️ [Background Push] Error enviando template a %s: %v\n", r.Number, err)
							continue
						}

						// Intentar actualizar el mensaje con el wamid retornado
						if msgID, ok := messageIDs[r.Number]; ok {
							if err := dbConn.Model(&models.Message{}).Where("id = ?", msgID).Update("channel_message_id", wamid).Error; err != nil {
								fmt.Printf("⚠️ [Background Push] Error actualizando channel_message_id para mensaje %d: %v\n", msgID, err)
							}
						}
					}
				}
			}()
		}
		wg.Wait()
		fmt.Printf("✅ [Background Push] Envío masivo finalizado para push ID: %d\n", pushID)
	}(recipients, messageIDs, integration, template)

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

	// create message record
	caseIDParam := int64(caseChannelIntegration.CaseID)

	var wabaID string
	if caseChannelIntegration.MetaWabaID != nil && *caseChannelIntegration.MetaWabaID != "" {
		wabaID = *caseChannelIntegration.MetaWabaID
	} else {
		wabaID, _ = GetSettingTextValue("WAB_ID")
	}

	bodyText, err := GetTemplateBodyFromMeta(*&template.TemplateName, wabaID, *caseChannelIntegration.AccessToken)
	if err != nil {
		bodyText = "Apertura mediante template"
	}

	templateIDUint := uint(templateID)
	newMessage := models.Message{
		CaseID:      uint(caseIDParam),
		SenderType:  "agent",
		MessageType: "text",
		TextContent: bodyText,
		MessageRead: true,
		TemplateID:  &templateIDUint,
	}

	if err := config.DB.Create(&newMessage).Error; err != nil {
		return nil, err
	}

	wamid, err := utils.SendTemplateMessageWithID(
		template.TemplateUrlWebhook,
		*caseChannelIntegration.AppIdentifier,
		*caseChannelIntegration.AccessToken,
		template.TemplateName,
		template.Language,
		*caseChannelIntegration.SenderID,
	)

	if err != nil {
		return nil, fmt.Errorf("error enviando template: %w", err)
	}

	// Actualizar el mensaje con el ID de Meta
	if err := config.DB.Model(&models.Message{}).Where("id = ?", newMessage.ID).Update("channel_message_id", wamid).Error; err != nil {
		fmt.Printf("⚠️ Error actualizando channel_message_id para mensaje %d: %v\n", newMessage.ID, err)
	}

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
		var existingCase models.Case
		var caseIDToUse uint

		err := tx.Debug().Where("channel_integration_id = ? AND sender_id = ? AND status = ?", request.ChannelIntegrationID, request.ContactPhone, "open").First(&existingCase).Error
		if err == nil {
			caseIDToUse = existingCase.ID
		} else if err == gorm.ErrRecordNotFound {

			// Create new case

			number := strconv.FormatUint(uint64(integration.ChannelID), 10)

			finalAgentID := uint(request.AgentID)
			agentAssigned := false

			// Intento 1: Asignación directa en ChannelAgentClient
			if request.ClientID != nil {
				var channelAgentClient models.ChannelAgentClient
				errAgent := tx.Where("department_id = ? AND client_id = ?", *integration.DepartmentID, *request.ClientID).First(&channelAgentClient).Error
				if errAgent == nil {
					finalAgentID = uint(channelAgentClient.AgentID)
					agentAssigned = true
				}
			}

			// Intento 2: Historial del cliente en el canal (último caso atendido)
			if !agentAssigned {
				var lastCase models.Case
				errLastCase := tx.Where("channel_integration_id = ? AND sender_id = ?", request.ChannelIntegrationID, request.ContactPhone).Order("id desc").First(&lastCase).Error
				if errLastCase == nil {
					// Verificar si el agente asignado a ese caso previo sigue activo
					var agent models.User
					errUser := tx.Where("id = ?", lastCase.AgentID).First(&agent).Error
					if errUser == nil && agent.IsActive {
						finalAgentID = lastCase.AgentID
					}
				}
			}

			newCase := models.Case{
				SenderId:             request.ContactPhone,
				ChannelID:            number,
				CompanyID:            integration.CompanyID,
				ChannelIntegrationID: integration.ChannelIntegrationID,
				IsNonCommercial:      integration.IsNonCommercial,
				DepartmentID:         *integration.DepartmentID,
				ClientID:             request.ClientID,
				AgentID:              finalAgentID,
				Status:               "open",
			}

			if err := tx.Create(&newCase).Error; err != nil {
				return err
			}
			caseIDToUse = newCase.ID

		} else {
			return err
		}

		var wabaID string
		if integration.MetaWabaID != nil && *integration.MetaWabaID != "" {
			wabaID = *integration.MetaWabaID
		} else {
			wabaID, _ = GetSettingTextValue("WAB_ID")
		}

		bodyText, err := GetTemplateBodyFromMeta(*template.TemplateName, wabaID, *template.AccessToken)
		if err != nil {
			bodyText = "Apertura mediante template"
		}

		var tempIDUint *uint
		if template.TemplateID != nil {
			val := uint(*template.TemplateID)
			tempIDUint = &val
		}

		newMessage := models.Message{
			CaseID:      caseIDToUse,
			SenderType:  "agent",
			MessageType: "text",
			TextContent: bodyText,
			MessageRead: true,
			TemplateID:  tempIDUint,
		}

		if err := tx.Create(&newMessage).Error; err != nil {
			return err
		}

		caseID = int64(caseIDToUse)

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

func GetTemplateBodyFromMeta(templateName string, wabaID string, accessToken string) (string, error) {
	if wabaID == "" {
		var err error
		wabaID, err = GetSettingTextValue("WAB_ID")
		if err != nil {
			return "", err
		}
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
