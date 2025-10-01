package repository

import (
	"fmt"
	"harmony_api/config"
	"harmony_api/models"

	"gorm.io/gorm"
)

type CampaignPushingRepository struct{}

func NewCampaignPushingRepository() *CampaignPushingRepository {
	return &CampaignPushingRepository{}
}

// CreateWhatsappPush guarda el encabezado y los leads en una sola transacción
func (r *CampaignPushingRepository) CreateWhatsappPush(data *models.CampaignWhatsappPushRequest) (int64, error) {
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

		// 2. Guardar leads si existen
		if len(data.Leads) > 0 {
			var leads []models.CampaignWhatsappPushLead
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

				lead := models.CampaignWhatsappPushLead{
					PushID:      pushID,
					PhoneNumber: l.PhoneNumber,
					ClientID:    clientID,
					CaseID:      nil,
					FullName:    l.FullName,
					MessageSent: false,
				}

				leads = append(leads, lead)
			}

			if err := tx.Create(&leads).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return pushID, nil
}

// package repository

// import (
// 	"fmt"
// 	"harmony_api/config"
// 	"harmony_api/models"
// 	"time"

// 	"gorm.io/gorm"
// )

// type CampaignPushingRepository struct{}

// func NewCampaignPushingRepository() *CampaignPushingRepository {
// 	return &CampaignPushingRepository{}
// }

// // CreateWhatsappPush guarda el encabezado y los leads
// func (r *CampaignPushingRepository) CreateWhatsappPush(data *models.CampaignWhatsappPushRequest) (int64, error) {
// 	db := config.DB
// 	var pushID int64

// 	err := db.Transaction(func(tx *gorm.DB) error {
// 		// Guardar encabezado
// 		header := models.CampaignWhatsappPush{
// 			CampaignID:    data.CampaignID,
// 			Description:   data.Description,
// 			TemplateID:    data.TemplateID,
// 			FunnelStageID: data.FunnelStageID,
// 			ChangedBy:     data.ChangedBy,
// 		}

// 		// Search channel ID by template ID
// 		var template models.CompanyChannelTemplateView

// 		if err := tx.Where("template_id = ?", data.TemplateID).First(&template).Error; err != nil {
// 			return err
// 		}

// 		header.CreatedAt = time.Now()

// 		if err := tx.Create(&header).Error; err != nil {
// 			return err
// 		}

// 		// Obtener ID generado
// 		pushID = header.ID

// 		//
// 		// Get funnel repository
// 		//
// 		campaignRepo := NewCampaignRepository()
// 		campaign, err := campaignRepo.GetByID(uint(data.CampaignID))

// 		if err != nil {
// 			return fmt.Errorf("error al obtener la campaña: %w", err)
// 		}

// 		if campaign == nil {
// 			return fmt.Errorf("campaña no encontrada con ID: %d", data.CampaignID)
// 		}

// 		// Guardar leads
// 		if len(data.Leads) > 0 {
// 			var leads []models.CampaignWhatsappPushLead
// 			for _, l := range data.Leads {

// 				var clienteChannel models.VWClientSocialAccount
// 				var clientID *int64
// 				var caseID *int64

// 				// Buscar cliente por canal + número
// 				if err := config.DB.
// 					Where("channel_id = ? AND social_external_id = ?", template.ChannelID, l.PhoneNumber).
// 					First(&clienteChannel).Error; err == nil {

// 					id := int64(clienteChannel.ClientID) // convertir a int64
// 					clientID = &id
// 				} else if err != gorm.ErrRecordNotFound {
// 					fmt.Println("Error al buscar el cliente:", err)
// 				}

// 				// Verificar si existe un caso abierto para ese cliente en esa campaña
// 				var caseRecord models.Case
// 				if clientID != nil {
// 					if err := config.DB.
// 						Where("client_id = ? AND company_id = ? AND campaign_id = ? AND status = ?", *clientID, template.CompanyID, data.CampaignID, "open").
// 						First(&caseRecord).Error; err == nil {
// 						cid := int64(caseRecord.ID)
// 						caseID = &cid
// 					} else if err != gorm.ErrRecordNotFound {
// 						fmt.Println("Error al buscar el caso:", err)
// 					}
// 				}

// 				if caseID == nil {
// 					// Si no existe un caso abierto, crear uno nuevo
// 					newCase := models.Case{
// 						SenderId:   l.PhoneNumber,
// 						ChannelID:  string(template.ChannelID),
// 						CompanyID:  uint(template.CompanyID),
// 						CampaignID: uint(data.CampaignID),
// 						FunnelID:   uint(*campaign.FunnelID),
// 						Status:     "open",
// 						ClientID:   nil,
// 						StartedAt:  time.Now(),
// 					}

// 					if header.FunnelStageID != nil {
// 						newCase.CurrentStageID = header.FunnelStageID
// 					}

// 					if clientID != nil {
// 						cid := uint(*clientID)
// 						newCase.ClientID = &cid
// 					}

// 					if err := config.DB.Debug().Create(&newCase).Error; err != nil {
// 						return fmt.Errorf("error al crear el caso: %w", err)
// 					}

// 					cid := int64(newCase.ID)

// 					caseID = &cid
// 					fmt.Println("✅ Nuevo caso creado para lead:", l.PhoneNumber)
// 				} else {
// 					fmt.Println("📌 Caso existente encontrado para lead:", l.PhoneNumber)
// 				}

// 				// Si no existe un caso abierto, crear uno nuevo

// 				if header.FunnelStageID != nil {

// 					messageRepo := MessageRepository{}

// 					currentStage, err := messageRepo.GetCurrentCaseFunnel(int(*caseID))

// 					if err != nil {
// 						return fmt.Errorf("error al obtener la etapa actual del caso: %w", err)
// 					}

// 					toStage := int(*header.FunnelStageID)

// 					note := "Pushed via WhatsApp Campaign"

// 					newFunnelStage := models.CaseFunnel{
// 						CaseID:      int(*caseID),
// 						FunnelID:    int(*campaign.FunnelID),
// 						FromStageID: currentStage.CurrentStageID,
// 						ToStageID:   &toStage,
// 						Note:        &note,
// 						ChangedBy:   1, // Aquí deberías obtener el ID del usuario que hace el cambio
// 						ChangedAt:   time.Now(),
// 						Action:      "move",
// 					}

// 					if err := config.DB.Create(&newFunnelStage).Error; err != nil {
// 						return fmt.Errorf("error al crear el registro en case_funnel: %w", err)
// 					}

// 				}

// 				// Crear el lead asociado al push

// 				lead := models.CampaignWhatsappPushLead{
// 					PushID:      pushID,
// 					PhoneNumber: l.PhoneNumber,
// 					ClientID:    clientID, // ahora es *int64
// 					CaseID:      l.CaseID,
// 					MessageSent: l.MessageSent,
// 				}

// 				leads = append(leads, lead)
// 			}

// 			if err := tx.Create(&leads).Error; err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return 0, err
// 	}

// 	return pushID, nil
// }
