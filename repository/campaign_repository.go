package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type CampaignRepository struct {
}

func NewCampaignRepository() *CampaignRepository {
	return &CampaignRepository{}
}

// GetByCompany
func (r *CampaignRepository) GetByCompany(companyID uint) (*[]models.CampaignWithFunnel, error) {

	var campaign *[]models.CampaignWithFunnel
	err := config.DB.Where("company_id = ?", companyID).Find(&campaign).Error
	return campaign, err
}

// GetByID
func (r *CampaignRepository) GetByID(id uint) (map[string]interface{}, error) {

	var result map[string]interface{}
	err := config.DB.Debug().Model(&models.CampaignWithFunnel{}).Where("campaign_id = ?", id).First(&result).Error
	return result, err
}

// Create
func (r *CampaignRepository) Create(data *models.Campaign) error {
	return config.DB.Create(&data).Error
}

// Update
func (r *CampaignRepository) Update(data *models.Campaign) error {
	return config.DB.Debug().Save(&data).Error
}

// Delete
func (r *CampaignRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Campaign{}, id).Error
}
