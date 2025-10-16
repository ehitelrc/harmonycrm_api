package repository

import (
	"harmony_api/config"
	"harmony_api/models"
)

type CaseDashboardRepository struct {
}

func NewCaseDashboardRepository() *CaseDashboardRepository {
	return &CaseDashboardRepository{}
}

func (r *CaseDashboardRepository) refreshMaterializedView() error {
	return config.DB.Exec("REFRESH MATERIALIZED VIEW vw_case_dashboard_by_company").Error
}

// GetDashboardByCompany devuelve el dashboard de una compañía específica
func (r *CaseDashboardRepository) GetDashboardByCompany(companyID int64) (*models.CompanyDashboard, error) {

	if err := r.refreshMaterializedView(); err != nil {
		return nil, err
	}

	var dashboard models.CompanyDashboard

	if err := config.DB.Table("vw_case_dashboard_by_company").
		Where("company_id = ?", companyID).
		First(&dashboard).Error; err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// GetAllDashboards devuelve todos los dashboards (una fila por compañía)
func (r *CaseDashboardRepository) GetAllDashboards() ([]models.CompanyDashboard, error) {
	if err := r.refreshMaterializedView(); err != nil {
		return nil, err
	}

	var dashboards []models.CompanyDashboard
	if err := config.DB.Table("vw_case_dashboard_by_company").
		Find(&dashboards).Error; err != nil {
		return nil, err
	}
	return dashboards, nil
}
