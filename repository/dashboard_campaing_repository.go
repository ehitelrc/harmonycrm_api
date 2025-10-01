package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type DashboardCampaignRepository struct{}

func NewDashboardCampaignRepository() *DashboardCampaignRepository {
	return &DashboardCampaignRepository{}
}

func (r *DashboardCampaignRepository) GetCampaignSummaryByCompany(companyID uint) ([]models.VwDashboardCampaignPerCompany, error) {
	var rows []models.VwDashboardCampaignPerCompany
	err := config.DB.Where("company_id = ?", companyID).Find(&rows).Error
	return rows, err
}

func (r *DashboardCampaignRepository) GetGeneralDashboardByCompany(companyID uint) ([]models.VwDashboardGeneralByCompany, error) {
	var rows []models.VwDashboardGeneralByCompany
	err := config.DB.Where("company_id = ?", companyID).Find(&rows).Error
	return rows, err
}

func (r *DashboardCampaignRepository) GetCampaignFunnelSummary(campaignID int) ([]models.CampaignFunnelSummary, error) {
	var rows []models.CampaignFunnelSummary
	err := config.DB.Where("campaign_id = ?", campaignID).Find(&rows).Error
	return rows, err
}
