package repository

import (
	"errors"
	"fmt"
	"harmony_api/config"
	"harmony_api/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MessageRepository struct {
	// Aquí puedes definir campos necesarios para la conexión a la base de datos
}

func (r *MessageRepository) CreateMessage(message models.IncomingMessage) (*models.Message, error) {

	// Vamos a investigar si viene de un QR de contacto inicial

	prefix := "ccc||--FCH--||ccc"

	isQR := false

	var qrCompanyID, qrCampaignID, qrDepartmentID, qrUserID int

	if strings.HasPrefix(message.TextMessage, prefix) {
		encryptedPart := strings.TrimPrefix(message.TextMessage, prefix)
		encryptedPart = strings.TrimSpace(encryptedPart)

		parts := strings.Split(encryptedPart, "|")

		if len(parts) >= 4 {

			// Parte 0 es de la compañía
			// Parte 1 es el ID de la campaña
			// Parte 2 es el ID del usuario (si viene)
			qrCompanyID, _ = strconv.Atoi(parts[0])
			qrCampaignID, _ = strconv.Atoi(parts[1])
			qrUserID, _ = strconv.Atoi(parts[2])
			qrDepartmentID, _ = strconv.Atoi(parts[3])

			isQR = true

		}
	}

	var channnel models.VWChannel

	// Buscar canal por app_identifier
	if err := config.DB.
		Where("app_identifier = ?", message.RecipientID).
		First(&channnel).Error; err != nil {

		return nil, fmt.Errorf("canal no encontrado: %w", err)
	}

	// Buscar cliente por canal y sender_id
	var clienteChannel models.VWClientSocialAccount
	var clientID *uint
	hasClient := false

	if err := config.DB.
		Where("channel_id = ? AND social_external_id = ?", channnel.ChannelID, message.SenderID).
		First(&clienteChannel).Error; err == nil {
		clientID = &clienteChannel.ClientID
		hasClient = true
	} else if err != gorm.ErrRecordNotFound {
		fmt.Println("Error al buscar el cliente:", err)
	}

	// Buscar si ya existe un caso
	var cases models.Case
	tx := config.DB.
		Where("channel_id = ? AND sender_id = ? and status = ?", channnel.ChannelID, message.SenderID, "open").
		First(&cases)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}

	var caseID uint

	if tx.RowsAffected == 0 {
		// Crear nuevo caso
		newCase := models.Case{
			SenderId:             message.SenderID,
			ChannelID:            channnel.ChannelID,
			CompanyID:            channnel.CompanyID,
			ChannelIntegrationID: channnel.IntegrationID,
			IsNonCommercial:      channnel.IsNonCommercial,
			DepartmentID:         *channnel.DepartmentID,
			Status:               "open",
		}

		if hasClient {
			newCase.ClientID = clientID

			// Has channel agent assigned?
			var channelAgentClient models.ChannelAgentClient

			err := config.DB.Debug().Where("department_id = ? AND client_id = ?", channnel.DepartmentID, *clientID).First(&channelAgentClient).Error

			if err == nil {
				fmt.Println("🔔 Asignando agente del canal al caso:", channelAgentClient.AgentID)
				agentID := int(channelAgentClient.AgentID)
				newCase.AgentID = uint(agentID)
			}
		}

		// Si viene de QR, asignar company_id y campaign_id

		if isQR {

			newContact := models.QrLead{
				CompanyID:    qrCompanyID,
				CampaignID:   qrCampaignID,
				DepartmentID: &qrDepartmentID,
				UserID:       &qrUserID,
				ContactPhone: message.SenderID,
				Status:       "pending",
			}

			if hasClient {
				newContact.ClientID = clientID
			}

			if err := config.DB.Create(&newContact).Error; err != nil {
				return nil, fmt.Errorf("error al crear el contacto desde QR: %w", err)
			}

			newCase.DepartmentID = uint(qrDepartmentID)

			newCase.CompanyID = uint(qrCompanyID)

			// Get campaing funnel_id
			var campaign models.Campaign
			if err := config.DB.First(&campaign, qrCampaignID).Error; err != nil {
				return nil, fmt.Errorf("error al obtener la campaña: %w", err)
			}

			// Customer social network

			if campaign.FunnelID != nil {
				newCase.FunnelID = uint(*campaign.FunnelID)
			}

			if hasClient {

				// Verify if the client already has a social account for this channel

				currentClientChannel := models.ClientSocialAccount{}

				err := config.DB.Where("client_id = ? AND channel_id = ?", *clientID, channnel.ChannelID).
					First(&currentClientChannel).Error

				if err == nil {
					// Ya existe, no hacer nada
				} else if err != gorm.ErrRecordNotFound {

					channelStr := channnel.ChannelID
					clientChannelId, err := strconv.Atoi(channelStr) // convierte string → int

					if err != nil {
						return nil, fmt.Errorf("error al convertir channel_id a int: %w", err)
					}

					clientChannel := models.ClientSocialAccount{
						ClientID:   *clientID,
						ChannelID:  uint(clientChannelId),
						ExternalID: message.SenderID,
						Username:   message.FirstName,
						IsActive:   true,
					}

					if err := config.DB.Create(&clientChannel).Error; err != nil {
						return nil, fmt.Errorf("error al crear la cuenta social del cliente: %w", err)
					}
				}

			}

			newCase.CampaignID = uint(qrCampaignID)

			message.TextMessage = "Hola, soy " + message.FirstName + " Me he contactado a través del QR."

		}

		if err := config.DB.Create(&newCase).Error; err != nil {
			return nil, fmt.Errorf("error al crear el caso: %w", err)
		}
		caseID = newCase.ID
		fmt.Println("✅ Nuevo caso creado:", caseID)

	} else {
		caseID = cases.ID
		fmt.Println("📌 Caso existente encontrado:", caseID)
	}

	// Crear mensaje
	newMessage := models.Message{
		CaseID:        caseID,
		SenderType:    message.SenderType,
		MessageType:   message.MessageType,
		TextContent:   message.TextMessage,
		FileURL:       message.FileURL,
		MIMEType:      message.MIMEType,
		Base64Content: message.Base64Content,
		AgentID:       &cases.AgentID,
		MessageRead:   false,
	}

	if err := config.DB.Create(&newMessage).Error; err != nil {
		return nil, fmt.Errorf("error al crear el mensaje: %w", err)
	}

	fmt.Println("✉️  Nuevo mensaje creado:", newMessage.ID)
	return &newMessage, nil
}

// Get case by id
func (r *MessageRepository) GetCaseByID(id uint) (*models.Case, error) {
	var caseItem models.Case
	if err := config.DB.First(&caseItem, id).Error; err != nil {
		return nil, err
	}
	return &caseItem, nil
}

// Cases without agents assigned by company
func (r *MessageRepository) GetUnassignedCasesByCompanyID(companyID int) ([]models.CaseWithChannel, error) {
	var unassignedCases []models.CaseWithChannel
	err := config.DB.Where("company_id = ? AND agent_id IS NULL AND status = ?", companyID, "open").Find(&unassignedCases).Error
	return unassignedCases, err
}

// Cases without agents assigned by company and department
func (r *MessageRepository) GetUnassignedCasesByCompanyAndDepartmentID(companyID int, departmentID int) ([]models.CaseWithChannel, error) {
	var unassignedCases []models.CaseWithChannel

	err := config.DB.Where("company_id = ? AND department_id = ? AND agent_id IS NULL AND status = ?", companyID, departmentID, "open").Find(&unassignedCases).Error

	return unassignedCases, err
}

// Open cases by company and department
func (r *MessageRepository) GetOpenCasesByCompanyAndDepartmentID(companyID int, departmentID int) ([]models.CaseWithChannel, error) {
	var openCases []models.CaseWithChannel

	err := config.DB.Where("company_id = ? AND department_id = ? AND status = ?", companyID, departmentID, "open").Find(&openCases).Error

	return openCases, err
}

// Open cases by company and department with pagination
func (r *MessageRepository) GetOpenCasesByCompanyAndDepartmentIDPaged(
	companyID int,
	departmentID int,
	limit int,
	offset int,
) ([]models.CaseWithChannel, int64, error) {

	var cases []models.CaseWithChannel
	var total int64

	baseQuery := config.DB.
		Model(&models.CaseWithChannel{}).
		Where(
			"company_id = ? AND department_id = ? AND status = ?",
			companyID,
			departmentID,
			"open",
		)

	// total (sin paginar)
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// página actual
	if err := baseQuery.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&cases).Error; err != nil {
		return nil, 0, err
	}

	return cases, total, nil
}

func (r *MessageRepository) GetOpenCasesStatsByCompanyAndDepartment(
	companyID int,
	departmentID int,
) (total int64, assigned int64, unassigned int64, err error) {

	base := config.DB.
		Table("cases").
		Where(
			"company_id = ? AND department_id = ? AND status = ?",
			companyID,
			departmentID,
			"open",
		)

	// total
	if err = base.Count(&total).Error; err != nil {
		return
	}

	// asignados
	if err = base.
		Where("agent_id IS NOT NULL").
		Count(&assigned).Error; err != nil {
		return
	}

	// no asignados
	unassigned = total - assigned

	return
}

// Mark messages as read by case_id
func (r *MessageRepository) MarkMessagesAsReadByCaseID(caseID string) error {
	err := config.DB.Model(&models.Message{}).
		Where("case_id = ? AND message_read = ?", caseID, false).
		Update("message_read", true).Error
	return err
}

func (r *MessageRepository) GetActiveCasesByAgentID(agentID string) ([]models.CaseWithChannel, error) {
	var activeCases []models.CaseWithChannel
	err := config.DB.Where("agent_id = ? AND status = ?", agentID, "open").Find(&activeCases).Error
	return activeCases, err
}

func (r *MessageRepository) GetMessagesByCaseID(caseID string) ([]models.Message, error) {
	var messages []models.Message
	err := config.DB.Where("case_id = ?", caseID).Order("id ASC").Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) SendMessageToPlatform(message models.AgentMessage) error {

	// transform AgentMessage to Message

	// Agent case

	caseData := models.Case{}

	err := config.DB.Where("id = ?", message.CaseID).First(&caseData).Error

	if err != nil {
		return fmt.Errorf("error al obtener el caso: %w", err)
	}

	agentID := caseData.AgentID

	println("Agente -->", agentID)

	newMessage := models.Message{
		CaseID:        message.CaseID,
		SenderType:    message.SenderType,
		MessageType:   message.MessageType,
		TextContent:   message.TextMessage,
		Base64Content: message.Base64Content,
		MIMEType:      message.MIMEType,
		AgentID:       &agentID,
		HasError:      message.HasError,
		MessageError:  message.MessageError,
	}

	err = config.DB.Create(&newMessage).Error
	if err != nil {
		return fmt.Errorf("error al enviar el mensaje: %w", err)
	}

	return nil
}

func (r *MessageRepository) AssignCaseToClient(input models.AssignCaseInput) error {

	// Opción 1: actualizar solo una columna
	if err := config.DB.
		Model(&models.Case{}).
		Where("id = ?", input.CaseID).
		Update("client_id", input.ClientID). // solo toca client_id
		Error; err != nil {
		return fmt.Errorf("error al asignar el caso al cliente: %w", err)
	}

	return nil
}

func (r *MessageRepository) AddCaseNote(note models.CaseNote) error {
	if err := config.DB.Create(&note).Error; err != nil {
		return fmt.Errorf("error al agregar la nota del caso: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetCaseNotesByCaseID(caseID string) ([]models.CaseNoteView, error) {
	var notes []models.CaseNoteView
	err := config.DB.Where("case_id = ?", caseID).Find(&notes).Error
	return notes, err
}

// Sugerencia: cambia la firma para recibir changedBy.
func (r *MessageRepository) AssignCaseToCampaign(caseID int, campaignID int, changedBy int) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1) Obtener funnel_id de la campaña
		var campaign struct {
			ID       int  `gorm:"column:id"`
			FunnelID *int `gorm:"column:funnel_id"` // por si la columna es nullable
		}
		if err := tx.
			Table("campaigns").
			Select("id, funnel_id").
			Where("id = ?", campaignID).
			Take(&campaign).Error; err != nil {
			return fmt.Errorf("no se pudo obtener la campaña %d: %w", campaignID, err)
		}
		if campaign.FunnelID == nil {
			return fmt.Errorf("la campaña %d no tiene funnel asignado", campaignID)
		}

		// 2) Actualizar campaign_id del caso
		if err := tx.Model(&models.Case{}).
			Where("id = ?", caseID).
			Updates(map[string]interface{}{
				"campaign_id": campaignID,
				"funnel_id":   *campaign.FunnelID,
			}).Error; err != nil {
			return fmt.Errorf("error al asignar el caso %d a la campaña %d: %w", caseID, campaignID, err)
		}

		// 3) Insertar log en case_funnel (acción 'assign')
		// entry := models.CaseFunnel{
		// 	CaseID:      caseID,
		// 	FunnelID:    *campaign.FunnelID,
		// 	FromStageID: nil,
		// 	ToStageID:   nil,
		// 	Note:        nil,
		// 	ChangedBy:   changedBy,
		// 	Action:      "assign",
		// 	// ChangedAt: lo pone la DB (DEFAULT now())
		// }
		// if err := tx.Create(&entry).Error; err != nil {
		// 	return fmt.Errorf("no se pudo crear el log case_funnel (assign): %w", err)
		// }

		return nil
	})
}

func (r *MessageRepository) AssignCaseToDepartment(caseID int, departmentID int, changedBy int) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1) Actualizar el caso con el nuevo department_id
		if err := tx.Model(&models.Case{}).
			Where("id = ?", caseID).
			Update("department_id", departmentID).Error; err != nil {
			return fmt.Errorf("error al asignar el caso %d al departamento %d: %w", caseID, departmentID, err)
		}

		return nil
	})
}

func (r *MessageRepository) AssignCaseToAgent(caseID int, agentID int, changedBy int, departmentID int) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {

		// 1) Actualizar SIEMPRE el caso (agente y departamento)
		if err := tx.Model(&models.Case{}).
			Where("id = ?", caseID).
			Updates(map[string]interface{}{
				"department_id": departmentID,
				"agent_id":      agentID,
			}).Error; err != nil {
			return fmt.Errorf("error al asignar el caso %d al agente %d: %w", caseID, agentID, err)
		}

		// 2) Cargar datos actualizados del caso
		var caseData models.Case
		if err := tx.Where("id = ?", caseID).First(&caseData).Error; err != nil {
			return fmt.Errorf("error al obtener el caso %d: %w", caseID, err)
		}

		// 3) Solo si hay cliente asociado, sincronizar channel_agent_clients
		if caseData.ClientID != nil {

			// Eliminar asignaciones previas
			if err := tx.
				Where("channel_id = ? AND client_id = ?", caseData.ChannelID, caseData.ClientID).
				Delete(&models.ChannelAgentClient{}).Error; err != nil {
				return fmt.Errorf("error al eliminar asignaciones previas: %w", err)
			}

			// Convertir channel_id
			channelIDInt, err := strconv.ParseInt(caseData.ChannelID, 10, 64)
			if err != nil {
				return fmt.Errorf("channel_id inválido, no se pudo convertir a int64: %v", err)
			}

			// Crear nueva relación
			channelAgentClient := models.ChannelAgentClient{
				ChannelID:    channelIDInt,
				AgentID:      int64(agentID),
				DepartmentID: int64(departmentID),
				ClientID:     int64(*caseData.ClientID),
			}

			if err := tx.Create(&channelAgentClient).Error; err != nil {
				return fmt.Errorf("error al crear el nuevo registro en channel_agent_clients: %w", err)
			}
		}

		// Si todo salió bien, se hace COMMIT automático
		return nil
	})
}

// GetCurrentCaseFunnel
func (r *MessageRepository) GetCurrentCaseFunnel(caseID int) (models.VWCaseCurrentStage, error) {
	var caseFunnel models.VWCaseCurrentStage
	err := config.DB.Where("case_id = ?", caseID).
		Order("last_changed_by DESC").
		First(&caseFunnel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No hay registros -> retornamos objeto vacío y nil error
		return models.VWCaseCurrentStage{}, nil
	}

	return caseFunnel, err
}

func (r *MessageRepository) SetCaseFunnelStage(payload models.CaseFunnel) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1) Actualizar el case con el nuevo funnel_id (si es necesario)
		if err := tx.Model(&models.Case{}).
			Where("id = ?", payload.CaseID).
			Update("funnel_id", payload.FunnelID).Error; err != nil {
			return fmt.Errorf("error al actualizar el funnel del caso %d: %w", payload.CaseID, err)
		}

		if err := tx.Create(&payload).Error; err != nil {
			return fmt.Errorf("no se pudo crear el log case_funnel (move): %w", err)
		}

		// Update the current_stage_id in the cases table
		if err := tx.Model(&models.Case{}).
			Where("id = ?", payload.CaseID).
			Update("current_stage_id", payload.ToStageID).Error; err != nil {
			return fmt.Errorf("error al actualizar el current_stage_id del caso %d: %w", payload.CaseID, err)
		}

		return nil
	})
}

func (r *MessageRepository) CloseCase(request models.CaseCloseRequest) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1) Actualizar el case con el nuevo estado 'closed'
		if err := tx.Model(&models.Case{}).
			Where("id = ?", request.CaseID).
			Updates(map[string]interface{}{
				"closed_at": gorm.Expr("NOW()"),
				"status":    "closed",
			}).Error; err != nil {
			return fmt.Errorf("error al cerrar el caso %d: %w", request.CaseID, err)
		}

		// 2) Insertar log en case_funnel (acción 'close')
		var funnelID *int = nil
		if request.FunnelID != nil && *request.FunnelID > 0 {
			funnelID = request.FunnelID
		}

		entry := models.CaseFunnel{
			CaseID:      request.CaseID,
			FunnelID:    funnelID, // 👈 ya es *int
			FromStageID: nil,
			ToStageID:   nil,
			Note:        &request.Note,
			ChangedBy:   request.ClosedBy,
			Action:      "close",
		}
		if err := tx.Create(&entry).Error; err != nil {
			return fmt.Errorf("no se pudo crear el log case_funnel (close): %w", err)
		}

		return nil
	})
}

func (r *MessageRepository) CancelCase(caseID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1) Actualizar el caso a estado 'cancelled'
		if err := tx.Model(&models.Case{}).
			Where("id = ?", caseID).
			Update("status", "cancelled").Error; err != nil {
			return fmt.Errorf("error al cancelar el caso %d: %w", caseID, err)
		}

		return nil
	})
}

// GetCaseGeneralInformation
func (r *MessageRepository) GetCaseGeneralInformation(companyID, campaignID uint) ([]models.VWCaseGeneralInformation, error) {
	var cases []models.VWCaseGeneralInformation
	err := config.DB.
		Where("company_id = ? AND campaign_id = ?", companyID, campaignID).
		Find(&cases).Error

	if err != nil {
		return nil, err
	}

	// aseguramos que aunque no haya resultados, devuelva slice vacío y no nil
	if cases == nil {
		cases = []models.VWCaseGeneralInformation{}
	}

	return cases, nil
}

// GetCaseGeneralInformationByCompanyCampaignAgent
func (r *MessageRepository) GetCaseGeneralInformationByCompanyCampaignAgent(companyID, campaignID, agentID uint, channelIntegrationID uint) ([]models.VWCaseGeneralInformation, error) {
	var cases []models.VWCaseGeneralInformation
	err := config.DB.
		Where("company_id = ? AND campaign_id = ? AND agent_id = ? AND channel_integration_id = ?", companyID, campaignID, agentID, channelIntegrationID).
		Find(&cases).Error

	if err != nil {
		return nil, err
	}

	// aseguramos que aunque no haya resultados, devuelva slice vacío y no nil
	if cases == nil {
		cases = []models.VWCaseGeneralInformation{}
	}

	return cases, nil
}

// vw_case_general_information by company_id, campaign_id and agent_id
func (r *MessageRepository) GetCaseGeneralInformationByAgent(companyID, campaignID, agentID uint) ([]models.VWCaseGeneralInformation, error) {
	var cases []models.VWCaseGeneralInformation
	err := config.DB.
		Where("company_id = ? AND campaign_id = ? AND agent_id = ?", companyID, campaignID, agentID).
		Find(&cases).Error

	if err != nil {
		return nil, err
	}

	// aseguramos que aunque no haya resultados, devuelva slice vacío y no nil
	if cases == nil {
		cases = []models.VWCaseGeneralInformation{}
	}

	return cases, nil
}

func (r *MessageRepository) GetMessageControl(messageID string) (bool, error) {
	var record models.WhatsAppMessageControl

	// 1️⃣ Buscar si existe
	err := config.DB.
		Where("ws_message__id = ?", messageID).
		First(&record).Error

	if err == nil {
		// ✔ YA existe
		return true, nil
	}

	// Si el error es distinto a "not found", es un error real
	if err != gorm.ErrRecordNotFound {
		return false, err
	}

	// 2️⃣ NO existe → Insertamos el control
	newRecord := models.WhatsAppMessageControl{
		WSMessageID: messageID,
	}

	if err := config.DB.Create(&newRecord).Error; err != nil {
		return false, err
	}

	// ✔ NO existía, ya quedó insertado
	return false, nil
}

func (r *MessageRepository) SaveOutgoingMessageStatus(
	caseID int,
	text string,
	wamid string,
	apiResponse string,
	err error,
) error {

	// status := "sent"
	// errMsg := ""

	// if err != nil {
	// 	status = "failed"
	// 	errMsg = err.Error()
	// }

	// query := `
	//     INSERT INTO outgoing_messages_log(case_id, message_text, wamid, api_response, status, error_message, created_at)
	//     VALUES ($1, $2, $3, $4, $5, $6, NOW())
	// `
	// _, execErr := config.DB.Exec(
	// 	query,
	// 	caseID,
	// 	text,
	// 	wamid,
	// 	apiResponse,
	// 	status,
	// 	errMsg,
	// )

	return nil
}

func (r *MessageRepository) UpdateMediaContent(
	messageID uint,
	base64 string,
	mime string,
) error {
	return config.DB.Model(&models.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"base64_content": base64,
			"mime_type":      mime,
		}).Error
}

func (r *MessageRepository) IsMessengerWindowOpen(caseID uint) (bool, error) {
	var lastInbound struct {
		CreatedAt time.Time
	}

	err := config.DB.
		Model(&models.Message{}).
		Select("created_at").
		Where("case_id = ?", caseID).
		Where("sender_type = ?", "client").
		Order("created_at DESC").
		Limit(1).
		Scan(&lastInbound).Error

	if err != nil {
		return false, err
	}

	// Si no hay mensajes entrantes → ventana cerrada
	if lastInbound.CreatedAt.IsZero() {
		return false, nil
	}

	timer := time.Since(lastInbound.CreatedAt)
	fmt.Printf("⏱ Tiempo desde el último mensaje entrante: %v\n", timer)

	// Ventana válida = ahora < last_inbound + 24h
	return time.Since(lastInbound.CreatedAt) <= 24*time.Hour, nil
}
