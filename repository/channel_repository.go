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

func (r *ChannelRepository) GetWhatsappTemplatesByChannelID(channelId uint) ([]models.ChannelWhatsAppTemplate, error) {
	var templates []models.ChannelWhatsAppTemplate
	if err := config.DB.Where("channel_integration_id = ?", channelId).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *ChannelRepository) UpdateWhatsappTemplate(template *models.ChannelWhatsAppTemplate) error {
	return config.DB.Save(template).Error
}

func (r *ChannelRepository) GetWhatsappTemplatesByIntegrationID(integrationId uint) ([]models.ChannelWhatsAppTemplate, error) {
	var templates []models.ChannelWhatsAppTemplate
	if err := config.DB.Debug().Where("channel_integration_id = ?", integrationId).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *ChannelRepository) DeleteWhatsappTemplate(template_id uint) error {
	return config.DB.Delete(&models.ChannelWhatsAppTemplate{}, template_id).Error
}

func (r *ChannelRepository) GetChannelWhatsappIntegrationsByCompanyID(companyId uint) ([]models.ViewChannelIntegration, error) {
	var integrations []models.ViewChannelIntegration
	if err := config.DB.Where("company_id = ? and channel_code = ?", companyId, "whatsapp").Find(&integrations).Error; err != nil {
		return nil, err
	}
	return integrations, nil
}

// By department
func (r *ChannelRepository) GetWhatsappTemplatesByDepartmentID(departmentId uint) ([]models.VwChannelWhatsAppTemplate, error) {
	var templates []models.VwChannelWhatsAppTemplate
	if err := config.DB.Where("department_id = ?", departmentId).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *ChannelRepository) GetChannelWhatsappIntegrationsByDepartmentID(departmentId uint) ([]models.ViewChannelIntegration, error) {
	var integrations []models.ViewChannelIntegration
	if err := config.DB.Where("department_id = ? and channel_code = ?", departmentId, "whatsapp").Find(&integrations).Error; err != nil {
		return nil, err
	}
	return integrations, nil
}

// By channel_integration_id

func (r *ChannelRepository) GetChannelIntegrationByID(integration_id uint) ([]models.ViewChannelIntegration, error) {
	var templates []models.ViewChannelIntegration
	if err := config.DB.Where("channel_integration_id  = ?", integration_id).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// By company_id and channel_id
func (r *ChannelRepository) GetChannelIntegrationsByCompanyAndChannelID(company_id uint, channel_id uint) ([]models.ViewChannelIntegration, error) {
	var integrations []models.ViewChannelIntegration
	if err := config.DB.Where("company_id = ? and channel_id = ?", company_id, channel_id).Find(&integrations).Error; err != nil {
		return nil, err
	}
	return integrations, nil
}

// By app_indentifier
func (r *ChannelRepository) GetChannelIntegrationByAppIdentifier(app_identifier string) (*models.ViewChannelIntegration, error) {
	var integration models.ViewChannelIntegration
	if err := config.DB.Where("app_identifier = ?", app_identifier).First(&integration).Error; err != nil {
		return nil, err
	}
	return &integration, nil
}

func (r *ChannelRepository) AddIntegrationToChannel(integration *models.ChannelIntegration) error {
	return config.DB.Create(integration).Error
}

func (r *ChannelRepository) UpdateChannelIntegration(integration *models.ChannelIntegration) error {
	return config.DB.Save(integration).Error
}
