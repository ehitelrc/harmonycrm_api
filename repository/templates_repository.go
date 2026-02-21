package repository

import (
	"errors"
	"harmony_api/config"
	"harmony_api/models"
)

type TemplateRepository struct{}

func NewTemplateRepository() *TemplateRepository {
	return &TemplateRepository{}
}

func (r *TemplateRepository) GetAll(channelID *uint) ([]models.MessageTemplate, error) {
	var templates []models.MessageTemplate

	q := config.DB.
		Select(`message_templates.*,
			(SELECT COUNT(*) FROM integration_templates it
			 WHERE it.template_id = message_templates.id) AS linked_count`).
		Model(&models.MessageTemplate{})

	if channelID != nil {
		q = q.Debug().Where("message_templates.channel_id = ?", *channelID)
	}

	err := q.Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) GetByID(id uint) (*models.MessageTemplate, error) {
	var template models.MessageTemplate
	err := config.DB.Preload("Channel").First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TemplateRepository) Create(template *models.MessageTemplate) error {
	return config.DB.Create(template).Error
}

func (r *TemplateRepository) Update(id uint, payload map[string]interface{}) error {
	return config.DB.Model(&models.MessageTemplate{}).Where("id = ?", id).Updates(payload).Error
}

func (r *TemplateRepository) Delete(id uint) error {
	return config.DB.Delete(&models.MessageTemplate{}, id).Error
}

// ---- Integration-Template relations ----

func (r *TemplateRepository) GetByIntegrationID(integrationID uint) ([]models.ChannelTemplateIntegration, error) {
	var results []models.ChannelTemplateIntegration
	err := config.DB.
		Where("integration_id = ?", integrationID).
		Find(&results).Error
	return results, err
}

// GetIntegrationsForTemplate returns all available integrations for a given template,
// including the is_linked flag from vw_channel_template_integrations.
func (r *TemplateRepository) GetIntegrationsForTemplate(templateID uint) ([]models.ChannelTemplateIntegration, error) {
	var results []models.ChannelTemplateIntegration
	err := config.DB.
		Where("template_id = ?", templateID).
		Find(&results).Error
	return results, err
}

func (r *TemplateRepository) CreateIntegrationTemplate(integrationID, templateID uint) (*models.IntegrationTemplate, error) {
	var existing models.IntegrationTemplate
	result := config.DB.
		Where("integration_id = ? AND template_id = ?", integrationID, templateID).
		First(&existing)
	if result.Error == nil {
		return nil, errors.New("la plantilla ya está asignada a esta integración")
	}

	rel := models.IntegrationTemplate{
		IntegrationID: integrationID,
		TemplateID:    templateID,
	}
	if err := config.DB.Create(&rel).Error; err != nil {
		return nil, err
	}
	return &rel, nil
}

func (r *TemplateRepository) DeleteIntegrationTemplate(integrationID, templateID uint) error {
	result := config.DB.
		Where("integration_id = ? AND template_id = ?", integrationID, templateID).
		Delete(&models.IntegrationTemplate{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("relación no encontrada")
	}
	return nil
}
