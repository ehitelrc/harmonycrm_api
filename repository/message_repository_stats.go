package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

func (r *MessageRepository) GetCaseStats(
	companyID uint,
	departmentID uint,
) (*models.CaseStatsResponse, error) {

	stats := models.CaseStatsResponse{}

	// 1️⃣ Summary
	err := config.DB.Raw(`
		SELECT
			COUNT(*) AS total_cases,
			COUNT(*) FILTER (WHERE status = 'open') AS open_cases,
			COUNT(*) FILTER (WHERE status = 'closed') AS closed_cases,
			COUNT(*) FILTER (WHERE agent_id IS NULL) AS unassigned_cases,
			COUNT(*) FILTER (WHERE unread_count > 0) AS unread_cases
		FROM mv_cases_with_channels
		WHERE company_id = ? AND department_id = ?
	`, companyID, departmentID).Scan(&stats.Summary).Error

	if err != nil {
		return nil, err
	}

	// 2️⃣ By channel
	err = config.DB.Raw(`
		SELECT channel_name AS label, COUNT(*) AS total
		FROM mv_cases_with_channels
		WHERE company_id = ? AND department_id = ?
		GROUP BY channel_name
		ORDER BY total DESC
	`, companyID, departmentID).Scan(&stats.ByChannel).Error

	if err != nil {
		return nil, err
	}

	// 3️⃣ By agent
	err = config.DB.Raw(`
		SELECT COALESCE(agent_full_name, 'Sin asignar') AS label, COUNT(*) AS total
		FROM mv_cases_with_channels
		WHERE company_id = ? AND department_id = ?
		GROUP BY label
		ORDER BY total DESC
	`, companyID, departmentID).Scan(&stats.ByAgent).Error

	if err != nil {
		return nil, err
	}

	// 4️⃣ By funnel_stage
	err = config.DB.Raw(`
		SELECT COALESCE(funnel_stage::text, 'Sin etapa') AS label, COUNT(*) AS total
		FROM mv_cases_with_channels
		WHERE company_id = ? AND department_id = ?
		GROUP BY label
		ORDER BY label
	`, companyID, departmentID).Scan(&stats.ByFunnelStage).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}
