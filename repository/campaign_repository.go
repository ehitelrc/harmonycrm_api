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
func (r *CampaignRepository) GetByID(id uint) (*models.CampaignWithFunnel, error) {

	var result *models.CampaignWithFunnel
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

// GetReachReport retrieves sent/delivered/read/replied stats for a campaign
func (r *CampaignRepository) GetReachReport(campaignID uint) (*models.CampaignReachResponse, error) {
	db := config.DB

	var recipients []models.CampaignReachRecipient
	query := `
		SELECT 
			l.phone_number,
			COALESCE(l.full_name, '') as full_name,
			COALESCE(l.case_id, 0) as case_id,
			COALESCE((
				SELECT 
					CASE 
						WHEN m.status IS NULL OR m.status = '' THEN 'sent'
						ELSE m.status
					END
				FROM messages m 
				WHERE m.case_id = l.case_id AND m.sender_type = 'agent'
				ORDER BY m.id ASC LIMIT 1
			), 'pending') as status,
			EXISTS(
				SELECT 1 FROM messages m2 
				WHERE m2.case_id = l.case_id AND m2.sender_type = 'client'
			) as replied,
			p.created_at
		FROM campaign_whatsapp_push_leads l
		JOIN campaign_whatsapp_push p ON p.id = l.push_id
		WHERE p.campaign_id = ?
		ORDER BY l.id DESC
	`
	err := db.Raw(query, campaignID).Scan(&recipients).Error
	if err != nil {
		return nil, err
	}

	var summary models.CampaignReachSummary
	for _, rec := range recipients {
		summary.SentCount++
		if rec.Status == "read" {
			summary.ReadCount++
			summary.DeliveredCount++
		} else if rec.Status == "delivered" {
			summary.DeliveredCount++
		} else if rec.Status == "failed" {
			summary.FailedCount++
		}
		if rec.Replied {
			summary.RepliedCount++
		}
	}

	return &models.CampaignReachResponse{
		Summary:    summary,
		Recipients: recipients,
	}, nil
}
