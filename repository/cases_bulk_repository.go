package repository

import (
	"harmony_api/config"
	"harmony_api/dto"
	"time"
)

type CasesBulkRepository struct {
}

func NewCasesBulkRepository() *CasesBulkRepository {
	return &CasesBulkRepository{}
}

func (r *CasesBulkRepository) SearchCasesToBulkClose(companyID, departmentID uint, startDate, endDate time.Time) ([]dto.BulkCloseCaseItem, error) {
	var results []dto.BulkCloseCaseItem

	query := `
		SELECT 
			c.id as case_id,
			ch.code as channel_logo_url,
			ch.name as channel_integration_name,
			c.started_at,
			c.sender_id,
			cl.full_name as client_name,
			COALESCE(sent.count, 0) as messages_sent_count,
			COALESCE(recv.count, 0) as messages_received_count,
			last_msg.created_at as last_message_date,
			last_msg.sender_type as last_message_type,
			c.payment_found as payment_found,
			COALESCE(pr.count, 0) as payment_records_count
		FROM cases c
		LEFT JOIN channel_integrations ci ON c.channel_integration_id = ci.id
		LEFT JOIN channels ch ON ci.channel_id = ch.id
		LEFT JOIN clients cl ON c.client_id = cl.id
		LEFT JOIN (
			SELECT case_id, COUNT(*) as count FROM messages WHERE sender_type = 'agent' GROUP BY case_id
		) sent ON c.id = sent.case_id
		LEFT JOIN (
			SELECT case_id, COUNT(*) as count FROM messages WHERE sender_type = 'client' GROUP BY case_id
		) recv ON c.id = recv.case_id
		LEFT JOIN (
			SELECT DISTINCT ON (case_id) case_id, created_at, sender_type 
			FROM messages 
			ORDER BY case_id, created_at DESC
		) last_msg ON c.id = last_msg.case_id
		LEFT JOIN (
			SELECT case_id, COUNT(*) as count FROM receipt_results GROUP BY case_id
		) pr ON c.id = pr.case_id
		WHERE c.status = 'open' 
		  AND c.company_id = ? 
		  AND c.department_id = ? 
		  AND c.created_at BETWEEN ? AND ?
		LIMIT 1000
	`

	if err := config.DB.Debug().Raw(query, companyID, departmentID, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *CasesBulkRepository) ExecuteBulkClose(caseIDs []uint) error {
	return config.DB.Exec("UPDATE cases SET status = 'closed', closed_at = NOW(), closed_in_bulk = true WHERE id IN ?", caseIDs).Error
}
