package repository

import (
	"fmt"
	"harmony_api/config"
	"harmony_api/models"
	"harmony_api/utils"
	"harmony_api/ws"

	"gorm.io/gorm"
)

// Sen

type CampaignPushingRepository struct{}

func NewCampaignPushingRepository() *CampaignPushingRepository {
	return &CampaignPushingRepository{}
}

// CreateWhatsappPush guarda el encabezado y los leads en una sola transacción
func (r *CampaignPushingRepository) CreateWhatsappPush(data *models.CampaignWhatsappPushRequest, hub *ws.Hub) (int64, error) {
	db := config.DB
	var pushID int64

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Guardar encabezado
		header := models.CampaignWhatsappPush{
			CampaignID:  data.CampaignID,
			Description: data.Description,
			TemplateID:  data.TemplateID,
			ChangedBy:   data.ChangedBy,
		}

		if err := tx.Create(&header).Error; err != nil {
			return err
		}

		// Search channel ID by template ID
		var template models.CompanyChannelTemplateView

		if err := tx.Where("template_id = ?", data.TemplateID).First(&template).Error; err != nil {
			return err
		}

		// Recuperar ID generado
		pushID = header.ID

		var numbers []string

		// 2. Guardar leads si existen
		if len(data.Leads) > 0 {
			var leads []models.CampaignWhatsappPushLead
			var messages []models.Message

			var recipients []models.TemplateRecipient

			for _, l := range data.Leads {

				var clienteChannel models.VWClientSocialAccount
				var clientID *int64

				// Buscar cliente por canal + número
				if err := config.DB.
					Where("channel_id = ? AND social_external_id = ?", template.ChannelID, l.PhoneNumber).
					First(&clienteChannel).Error; err == nil {

					id := int64(clienteChannel.ClientID) // convertir a int64
					clientID = &id
				} else if err != gorm.ErrRecordNotFound {
					fmt.Println("Error al buscar el cliente:", err)
				}

				var caseId *int64 = nil

				if l.CaseID != nil {
					caseId = l.CaseID
				}

				println(caseId)

				lead := models.CampaignWhatsappPushLead{
					PushID:      pushID,
					PhoneNumber: l.PhoneNumber,
					ClientID:    clientID,
					CaseID:      caseId,
					FullName:    l.FullName,
					MessageSent: false,
				}

				leads = append(leads, lead)
				numbers = append(numbers, l.PhoneNumber)

				// Agregar a la lista de destinatarios del template
				recipients = append(recipients, models.TemplateRecipient{
					Number: l.PhoneNumber,
					CaseID: caseId,
				})

				if l.ManualStartingLead {
					messages = append(messages, models.Message{
						CaseID:      uint(*caseId),
						SenderType:  "client",
						MessageType: "text",
						TextContent: "Apertura mediante template",
					})

				}
			}

			if err := tx.Create(&leads).Error; err != nil {
				return err
			}

			if len(messages) > 0 {
				if err := tx.Create(&messages).Error; err != nil {
					return err
				}

			}

			utils.SendTemplateToMany(template.TemplateUrlWebhook, *template.AppIdentifier, *template.AccessToken, *template.TemplateName, *template.Language, recipients, hub)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return pushID, nil
}
