package repository

import (
	"fmt"
	"harmony_api/config"
	"harmony_api/models"

	"gorm.io/gorm"
)

type MessageRepository struct {
	// Aquí puedes definir campos necesarios para la conexión a la base de datos
}

func (r *MessageRepository) CreateMessage(message models.IncomingMessage) (*models.Message, error) {
	var channnel models.VWChannel

	// Buscar canal por app_identifier
	if err := config.DB.
		Where("app_identifier = ?", message.RecipientID).
		First(&channnel).Error; err != nil {
		defer config.CloseDB()
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
		Where("channel_id = ? AND sender_id = ?", channnel.ChannelID, message.SenderID).
		First(&cases)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}

	var caseID uint

	if tx.RowsAffected == 0 {
		// Crear nuevo caso
		newCase := models.Case{
			SenderId:  message.SenderID,
			ChannelID: channnel.ChannelID,
			CompanyID: channnel.CompanyID,
			Status:    "open",
		}
		if hasClient {
			newCase.ClientID = clientID
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
	}

	if err := config.DB.Create(&newMessage).Error; err != nil {
		return nil, fmt.Errorf("error al crear el mensaje: %w", err)
	}

	fmt.Println("✉️  Nuevo mensaje creado:", newMessage.ID)
	return &newMessage, nil
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

	newMessage := models.Message{
		CaseID:      message.CaseID,
		SenderType:  message.SenderType,
		MessageType: message.MessageType,
		TextContent: message.TextMessage,
	}

	err := config.DB.Create(&newMessage).Error
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
			Update("campaign_id", campaignID).Error; err != nil {
			return fmt.Errorf("error al asignar el caso %d a la campaña %d: %w", caseID, campaignID, err)
		}

		// 3) Insertar log en case_funnel (acción 'assign')
		entry := models.CaseFunnel{
			CaseID:      caseID,
			FunnelID:    *campaign.FunnelID,
			FromStageID: nil,
			ToStageID:   nil,
			Note:        nil,
			ChangedBy:   changedBy,
			Action:      "assign",
			// ChangedAt: lo pone la DB (DEFAULT now())
		}
		if err := tx.Create(&entry).Error; err != nil {
			return fmt.Errorf("no se pudo crear el log case_funnel (assign): %w", err)
		}

		return nil
	})
}

//GetCurrentCaseFunnel

func (r *MessageRepository) GetCurrentCaseFunnel(caseID int) (models.VWCaseCurrentStage, error) {
	var caseFunnel models.VWCaseCurrentStage
	err := config.DB.Where("case_id = ?", caseID).Order("last_changed_by DESC").First(&caseFunnel).Error
	return caseFunnel, err
}
