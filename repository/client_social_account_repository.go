package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type ClientSocialAccountRepository struct{}

func NewClientSocialAccountRepository() *ClientSocialAccountRepository {
	return &ClientSocialAccountRepository{}
}

// Listados básicos
func (r *ClientSocialAccountRepository) GetAll() ([]models.ClientSocialAccount, error) {
	var rows []models.ClientSocialAccount
	err := config.DB.Find(&rows).Error
	return rows, err
}

func (r *ClientSocialAccountRepository) GetByID(id uint) (*models.ClientSocialAccount, error) {
	var row models.ClientSocialAccount
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Filtros útiles
func (r *ClientSocialAccountRepository) GetByClient(clientID uint) ([]models.ClientSocialAccount, error) {
	var rows []models.ClientSocialAccount
	err := config.DB.Where("client_id = ?", clientID).Find(&rows).Error
	return rows, err
}

func (r *ClientSocialAccountRepository) GetByChannel(channelID uint) ([]models.ClientSocialAccount, error) {
	var rows []models.ClientSocialAccount
	err := config.DB.Where("channel_id = ?", channelID).Find(&rows).Error
	return rows, err
}

func (r *ClientSocialAccountRepository) GetByChannelAndExternal(channelID uint, externalID string) (*models.ClientSocialAccount, error) {
	var row models.ClientSocialAccount
	if err := config.DB.Where("channel_id = ? AND external_id = ?", channelID, externalID).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// CRUD
func (r *ClientSocialAccountRepository) Create(m *models.ClientSocialAccount) error {
	return config.DB.Create(m).Error
}

func (r *ClientSocialAccountRepository) Update(m *models.ClientSocialAccount) error {
	return config.DB.Save(m).Error
}

func (r *ClientSocialAccountRepository) Delete(id uint) error {
	return config.DB.Delete(&models.ClientSocialAccount{}, id).Error
}

// Activar/Desactivar (opcional)
func (r *ClientSocialAccountRepository) SetActive(id uint, active bool) error {
	return config.DB.Model(&models.ClientSocialAccount{}).
		Where("id = ?", id).
		Update("is_active", active).Error
}
