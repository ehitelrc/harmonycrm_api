package repository

import (
	"fmt"
	"harmony_api/config"
	"harmony_api/models"
	"strconv"

	"gorm.io/gorm"
)

type ClientRepository struct{}

func NewClientRepository() *ClientRepository { return &ClientRepository{} }

func (r *ClientRepository) GetAll() ([]models.Client, error) {
	var rows []models.Client
	err := config.DB.Find(&rows).Error
	return rows, err
}

func (r *ClientRepository) GetByID(id uint) (*models.Client, error) {
	var row models.Client
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ClientRepository) Create(m *models.Client) error {
	return config.DB.Create(m).Error
}

// Update recibe el objeto completo (incluye id)
func (r *ClientRepository) Update(m *models.Client) error {
	return config.DB.Save(m).Error
}

func (r *ClientRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Client{}, id).Error
}

// CreateLead crea un nuevo lead en la tabla vw_case_general_information
// CreateLead crea un nuevo lead y su caso asociado
func (r *ClientRepository) CreateLead(lead *models.LeadRequest) error {
	tx := config.DB.Begin()

	// 🔹 1. Verificar cliente
	var client models.Client
	if err := tx.First(&client, lead.ClientID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("cliente no encontrado: %w", err)
	}

	// 🔹 2. Verificar integración de canal
	var channelIntegration models.ChannelIntegration
	if err := tx.First(&channelIntegration, lead.ChannelIntegrationID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("integración de canal no encontrada: %w", err)
	}

	// 🔹 3. Verificar si ya existe un caso abierto
	var existingCase models.Case
	err := tx.
		Where("channel_id = ? AND client_id = ? AND status = ? AND campaign_id = ? AND agent_id = ?",
			lead.ChannelID, client.ID, "open", lead.CampaignID, lead.AgentID).
		First(&existingCase).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	// Si ya existe un caso abierto, no crear otro
	if err == nil {
		tx.Rollback()
		return nil
	}

	// Get campaign to obtain funnel ID
	var campaign models.Campaign
	if err := tx.First(&campaign, lead.CampaignID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("campaña no encontrada: %w", err)
	}

	// 🔹 4. Crear nuevo caso
	newCase := models.Case{
		ClientID:             &lead.ClientID,
		CompanyID:            lead.CompanyID,
		CampaignID:           lead.CampaignID,
		ChannelID:            strconv.Itoa(int(lead.ChannelID)),
		ChannelIntegrationID: lead.ChannelIntegrationID,
		AgentID:              lead.AgentID,
		IsNonCommercial:      channelIntegration.IsNonCommercial,
		Status:               "open",
		SenderId:             client.Phone,
		FunnelID:             *campaign.FunnelID,
		ManualStartingLead:   true,
	}

	if err := tx.Create(&newCase).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("no se pudo crear el caso: %w", err)
	}

	// 🔹 5. Insertar artículos (si existen)
	for _, item := range lead.Items {
		caseItem := models.CaseItem{
			CaseID:   int(newCase.ID),
			ItemID:   int(item.ItemID),
			Quantity: float64(item.Quantity),
			Price:    item.ItemPrice,
			Notes:    &item.Notes,
		}
		if err := tx.Create(&caseItem).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("no se pudo crear el artículo del caso: %w", err)
		}
	}

	// 🔹 6. Registrar mensaje inicial
	newMessage := models.Message{
		CaseID:      newCase.ID,
		SenderType:  "agent",
		MessageType: "text",
		TextContent: "No hay conversación iniciada con el cliente.",
	}

	if err := tx.Create(&newMessage).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("no se pudo crear el mensaje inicial: %w", err)
	}

	return tx.Commit().Error
}

// repository/client_repository.go
func (r *ClientRepository) GetClientsWithDuplicatePhones() ([]models.Client, error) {
	var rows []models.Client

	err := config.DB.
		Where(`
			phone IN (
				SELECT phone
				FROM clients
				WHERE phone IS NOT NULL
				  AND TRIM(phone) <> ''
				GROUP BY phone
				HAVING COUNT(*) > 1
			)
		`).
		Order("phone ASC, id ASC").
		Find(&rows).Error

	return rows, err
}

// repository/client_repository.go
func (r *ClientRepository) GetDuplicatePhonesDTO() ([]models.DuplicatePhoneGroupDTO, error) {

	type row struct {
		ID            uint
		ExternalID    *string
		FullName      *string
		Email         *string
		Phone         *string
		CreatedAt     string
		UpdatedAt     string
		CountryID     *uint
		ProvinceID    *uint
		CantonID      *uint
		DistrictID    *uint
		AddressDetail *string
		PostalCode    *string
		IsCitizen     bool
	}

	var rows []row

	query := `
		SELECT *
		FROM clients
		WHERE phone IN (
			SELECT phone
			FROM clients
			WHERE phone IS NOT NULL
			  AND TRIM(phone) <> ''
			GROUP BY phone
			HAVING COUNT(*) > 1
		)
		ORDER BY phone, id
	`

	if err := config.DB.Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}

	// Agrupar por teléfono
	groupMap := make(map[string]*models.DuplicatePhoneGroupDTO)

	for _, r := range rows {
		if r.Phone == nil {
			continue
		}

		grp, ok := groupMap[*r.Phone]
		if !ok {
			grp = &models.DuplicatePhoneGroupDTO{
				Phone:   *r.Phone,
				Clients: []models.ClientDTO{},
			}
			groupMap[*r.Phone] = grp
		}

		grp.Clients = append(grp.Clients, models.ClientDTO{
			ID:            r.ID,
			ExternalID:    r.ExternalID,
			FullName:      r.FullName,
			Email:         r.Email,
			Phone:         r.Phone,
			CreatedAt:     r.CreatedAt,
			UpdatedAt:     r.UpdatedAt,
			CountryID:     r.CountryID,
			ProvinceID:    r.ProvinceID,
			CantonID:      r.CantonID,
			DistrictID:    r.DistrictID,
			AddressDetail: r.AddressDetail,
			PostalCode:    r.PostalCode,
			IsCitizen:     r.IsCitizen,
		})
	}

	// Convertir map → slice + count
	out := make([]models.DuplicatePhoneGroupDTO, 0, len(groupMap))
	for _, g := range groupMap {
		g.Count = len(g.Clients)
		out = append(out, *g)
	}

	return out, nil
}

// repository/client_repository.go

func (r *ClientRepository) GetByExternalID(externalID string) ([]models.Client, error) {
	var rows []models.Client

	err := config.DB.
		Where("external_id = ?", externalID).
		Order("id ASC").
		Find(&rows).Error

	return rows, err
}
