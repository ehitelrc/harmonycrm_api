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
