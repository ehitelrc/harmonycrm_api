package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type TemplateReportRepository struct{}

func NewTemplateReportRepository() *TemplateReportRepository {
	return &TemplateReportRepository{}
}

func (r *TemplateReportRepository) GetTemplateReport(companyID int64) (*models.TemplateReportResponse, error) {
	var bulkSends []models.BulkTemplateSend
	var individualSends []models.IndividualTemplateSend

	// Query for bulk sends
	bulkSQL := `
		SELECT 
			p.id AS push_id,
			p.description AS description,
			COALESCE(u.full_name, 'Sistema') AS agent_name,
			t.template_name AS template_name,
			p.created_at AS created_at,
			COUNT(l.id) AS total_recipients,
			COUNT(l.id) FILTER (WHERE l.message_sent = true) AS successful_sends,
			COUNT(l.id) FILTER (WHERE l.message_sent = false) AS failed_sends
		FROM campaign_whatsapp_push p
		JOIN users u ON u.id = p.changed_by
		JOIN message_templates t ON t.id = p.template_id
		LEFT JOIN campaign_whatsapp_push_leads l ON l.push_id = p.id
		JOIN channel_integrations ci ON ci.id = p.channel_integration_id
		WHERE ci.company_id = ?
		GROUP BY p.id, p.description, u.full_name, t.template_name, p.created_at
		ORDER BY p.created_at DESC
	`
	if err := config.DB.Raw(bulkSQL, companyID).Scan(&bulkSends).Error; err != nil {
		return nil, err
	}

	// Query for individual sends
	indivSQL := `
		SELECT 
			m.id AS message_id,
			m.case_id AS case_id,
			COALESCE(u.full_name, 'Sistema') AS agent_name,
			t.template_name AS template_name,
			m.created_at AS created_at,
			c.sender_id AS client_phone,
			COALESCE(m.status, 'sent') AS status
		FROM messages m
		JOIN cases c ON c.id = m.case_id
		JOIN message_templates t ON t.id = m.template_id
		LEFT JOIN users u ON u.id = m.agent_id
		WHERE c.company_id = ?
		ORDER BY m.created_at DESC
	`
	if err := config.DB.Raw(indivSQL, companyID).Scan(&individualSends).Error; err != nil {
		return nil, err
	}

	return &models.TemplateReportResponse{
		BulkSends:       bulkSends,
		IndividualSends: individualSends,
	}, nil
}
