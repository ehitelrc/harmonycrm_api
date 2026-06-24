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

func (r *CaseDashboardRepository) refreshMaterializedWithDepartmentView() error {
	return config.DB.Exec("REFRESH MATERIALIZED VIEW vw_case_dashboard_by_company_with_department").Error
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

// By company and department
func (r *CaseDashboardRepository) GetDashboardByCompanyAndDepartment(companyID int64, departmentID int64) (*models.CompanyDashboard, error) {
	if err := r.refreshMaterializedWithDepartmentView(); err != nil {
		return nil, err
	}

	var dashboard models.CompanyDashboard
	if err := config.DB.Debug().
		Table("public.vw_case_dashboard_by_company_with_department").
		Where("company_id = ? AND department_id = ?", companyID, departmentID).
		First(&dashboard).Error; err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// By company and user

func (r *CaseDashboardRepository) GetDashboardByCompanyAndUser(companyID int64, userID int64) (*models.CompanyDashboard, error) {
	if err := r.refreshMaterializedView(); err != nil {
		return nil, err
	}

	var dashboard models.CompanyDashboard
	if err := config.DB.Table("vw_case_dashboard_by_company").
		Where("company_id = ? AND user_id = ?", companyID, userID).
		First(&dashboard).Error; err != nil {
		return nil, err
	}

	var agentDepartments models.AgentDepartmentAssignment

	if err := config.DB.Table("agent_department_assignments").
		Where("agent_id = ?", userID).
		First(&agentDepartments).Error; err != nil {
		return nil, err
	}

	// Filtrar los datos del dashboard según los departamentos asignados al agente

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

// GetCasesByStatus obtiene los casos detallados filtrados por compañía, departamento, estado y búsqueda
func (r *CaseDashboardRepository) GetCasesByStatus(companyID int64, departmentID *int64, status string, search string, page int, limit int) ([]models.CaseWithChannel, int64, error) {
	var cases []models.CaseWithChannel
	var total int64

	db := config.DB.Model(&models.CaseWithChannel{}).Where("company_id = ?", companyID)

	if departmentID != nil && *departmentID > 0 {
		db = db.Where("department_id = ?", *departmentID)
	}

	switch status {
	case "open":
		db = db.Where("status = ?", "open")
	case "closed":
		db = db.Where("status = ?", "closed")
	case "unanswered":
		db = db.Where("status IN ('open', 'in_progress') AND last_message_sender_type = 'client'")
	}

	if search != "" {
		searchQuery := "%" + search + "%"
		db = db.Where("client_name ILIKE ? OR sender_id ILIKE ? OR last_message_text ILIKE ? OR agent_full_name ILIKE ? OR CAST(case_id AS TEXT) LIKE ?", searchQuery, searchQuery, searchQuery, searchQuery, "%"+search+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&cases).Error; err != nil {
		return nil, 0, err
	}

	return cases, total, nil
}
