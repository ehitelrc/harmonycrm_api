package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type ChannelRepository struct{}

func (r *ChannelRepository) GetAllChannels() ([]models.Channel, error) {
	var channels []models.Channel
	if err := config.DB.Find(&channels).Error; err != nil {
		return nil, err
	}
	return channels, nil
}

func (r *ChannelRepository) GetChannelByID(id uint) (*models.Channel, error) {
	var channel models.Channel
	if err := config.DB.First(&channel, id).Error; err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *ChannelRepository) CreateChannel(channel *models.Channel) error {
	return config.DB.Create(channel).Error
}

func (r *ChannelRepository) UpdateChannel(channel *models.Channel) error {
	return config.DB.Save(channel).Error
}

func (r *ChannelRepository) DeleteChannel(id uint) error {
	return config.DB.Delete(&models.Channel{}, id).Error
}

func (r *ChannelRepository) GetChannerlByCaseID(caseId uint) (*models.VWCaseChannelIntegration, error) {
	var channelIntegration models.VWCaseChannelIntegration
	if err := config.DB.Where("case_id = ?", caseId).First(&channelIntegration).Error; err != nil {
		return nil, err
	}
	return &channelIntegration, nil
}

func (r *ChannelRepository) CreateWhatsappTemplate(template *models.ChannelWhatsAppTemplate) error {
	return config.DB.Create(template).Error
}

func (r *ChannelRepository) GetWhatsappTemplatesByCompanyID(companyId uint) ([]models.CompanyChannelTemplateView, error) {
	var templates []models.CompanyChannelTemplateView
	if err := config.DB.Where("company_id = ?", companyId).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *ChannelRepository) GetChannelWhatsappIntegrationsByCompanyID(companyId uint) ([]models.VWChannelIntegration, error) {
	var integrations []models.VWChannelIntegration
	if err := config.DB.Where("company_id = ? and channel_code = ?", companyId, "whatsapp").Find(&integrations).Error; err != nil {
		return nil, err
	}
	return integrations, nil
}

func (r *ChannelRepository) GetChannelIntegrationByID(integration_id uint) ([]models.ChannelWhatsAppTemplate, error) {
	var templates []models.ChannelWhatsAppTemplate
	if err := config.DB.Where("channel_integration = ?", integration_id).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}
