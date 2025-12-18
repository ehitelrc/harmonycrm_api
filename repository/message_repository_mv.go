package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

func (r *MessageRepository) GetOpenCasesMV(
	companyID uint,
	departmentID uint,
	limit int,
	offset int,
) ([]models.CaseWithChannelMV, error) {

	var cases []models.CaseWithChannelMV

	err := config.DB.
		Where("company_id = ? AND department_id = ? AND status = ?", companyID, departmentID, "open").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&cases).Error

	return cases, err
}

func (r *MessageRepository) RefreshCasesMV(concurrent bool) error {

	sql := "REFRESH MATERIALIZED VIEW mv_cases_with_channels"

	if concurrent {
		sql = "REFRESH MATERIALIZED VIEW CONCURRENTLY mv_cases_with_channels"
	}

	return config.DB.Exec(sql).Error
}

func (r *MessageRepository) RefreshCasesMVConcurrently() error {
	return r.RefreshCasesMV(true)
}

func (r *MessageRepository) RefreshCasesMVBlocking() error {
	return r.RefreshCasesMV(false)
}
