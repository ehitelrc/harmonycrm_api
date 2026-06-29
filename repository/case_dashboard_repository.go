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
func (r *CaseDashboardRepository) GetDashboardByCompany(companyID int64, startDate, endDate string) (*models.CompanyDashboard, error) {

	if startDate == "" && endDate == "" {
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

	// Dynamic query with date filter
	var dashboard models.CompanyDashboard
	sqlQuery := `
WITH channel_stats AS (
    SELECT cs_1.company_id,
        ch.id AS channel_id,
        ch.name AS channel_name,
        count(cs_1.id) FILTER (WHERE (cs_1.status = 'open'::text)) AS open_cases,
        count(cs_1.id) FILTER (WHERE (cs_1.status = 'closed'::text)) AS closed_cases
    FROM public.cases cs_1
    LEFT JOIN public.channels ch ON ((ch.id = cs_1.channel_id))
    WHERE cs_1.company_id = ? AND cs_1.created_at >= ? AND cs_1.created_at <= ?
    GROUP BY cs_1.company_id, ch.id, ch.name
), agent_stats AS (
    SELECT cs_1.company_id,
        u.id AS agent_id,
        u.full_name AS agent_name,
        count(cs_1.id) FILTER (WHERE (cs_1.status = 'open'::text)) AS open_cases,
        count(cs_1.id) FILTER (WHERE (cs_1.status = 'closed'::text)) AS closed_cases,
        round(avg((EXTRACT(epoch FROM (cs_1.closed_at - CASE WHEN cs_1.started_at IS NULL OR cs_1.started_at <= '1970-01-02'::timestamp THEN cs_1.created_at ELSE cs_1.started_at END)) / (3600)::numeric)), 2) AS avg_close_hours
    FROM public.cases cs_1
    LEFT JOIN public.users u ON ((u.id = cs_1.agent_id))
    WHERE (cs_1.agent_id IS NOT NULL) AND cs_1.company_id = ? AND cs_1.created_at >= ? AND cs_1.created_at <= ?
    GROUP BY cs_1.company_id, u.id, u.full_name
), oldest_open AS (
    SELECT cs_1.company_id,
        cs_1.id AS case_id,
        cl.full_name AS client_name,
        cl.phone AS client_phone,
        cs_1.created_at,
        (SELECT m.created_at FROM public.messages m WHERE m.case_id = cs_1.id ORDER BY m.id DESC LIMIT 1) AS last_message_at,
        (SELECT m.sender_type FROM public.messages m WHERE m.case_id = cs_1.id ORDER BY m.id DESC LIMIT 1) AS last_message_sender_type
    FROM public.cases cs_1
    LEFT JOIN public.clients cl ON ((cl.id = cs_1.client_id))
    WHERE (cs_1.status = 'open'::text) AND cs_1.company_id = ? AND cs_1.created_at >= ? AND cs_1.created_at <= ?
)
SELECT c.id AS company_id,
    c.name AS company_name,
    count(cs.id) AS total_cases,
    count(cs.id) FILTER (WHERE (cs.status = 'open'::text)) AS open_cases,
    count(cs.id) FILTER (WHERE (cs.status = 'closed'::text)) AS closed_cases,
    count(cs.id) FILTER (WHERE ((cs.status = 'closed'::text) AND (date(cs.closed_at) = CURRENT_DATE))) AS closed_today,
    count(cs.id) FILTER (WHERE ((cs.status = 'open'::text) AND (date(cs.created_at) = CURRENT_DATE))) AS opened_today,
    count(cs.id) FILTER (WHERE (cs.status = 'cancelled'::text)) AS cancelled_cases,
    count(cs.id) FILTER (WHERE (
        cs.status IN ('open', 'in_progress')
        AND (
            SELECT m.sender_type
            FROM public.messages m
            WHERE m.case_id = cs.id
            ORDER BY m.id DESC
            LIMIT 1
        ) = 'client'
    )) AS unanswered_cases,
    count(cs.id) FILTER (WHERE (cs.agent_id IS NULL)) AS unassigned_agents,
    count(cs.id) FILTER (WHERE (cs.client_id IS NULL)) AS unassigned_clients,
    round(avg((EXTRACT(epoch FROM (cs.closed_at - CASE WHEN cs.started_at IS NULL OR cs.started_at <= '1970-01-02'::timestamp THEN cs.created_at ELSE cs.started_at END)) / (3600)::numeric)), 2) AS avg_close_hours,
    ( SELECT json_agg(json_build_object('channel_id', ch.channel_id, 'channel_name', ch.channel_name, 'open_cases', ch.open_cases, 'closed_cases', ch.closed_cases))
      FROM channel_stats ch
      WHERE (ch.company_id = c.id)) AS cases_by_channel,
    ( SELECT json_agg(json_build_object('agent_id', a.agent_id, 'agent_name', a.agent_name, 'open_cases', a.open_cases, 'closed_cases', a.closed_cases, 'avg_close_hours', a.avg_close_hours))
      FROM agent_stats a
      WHERE (a.company_id = c.id)) AS cases_by_agent,
    ( SELECT json_agg(json_build_object('case_id', o.case_id, 'client_name', o.client_name, 'client_phone', o.client_phone, 'created_at', o.created_at, 'last_message_at', o.last_message_at, 'last_message_sender_type', o.last_message_sender_type) ORDER BY o.created_at)
      FROM oldest_open o
      WHERE (o.company_id = c.id)
      LIMIT 20) AS oldest_open_cases
FROM public.companies c
LEFT JOIN public.cases cs ON ((cs.company_id = c.id) AND cs.created_at >= ? AND cs.created_at <= ?)
WHERE c.id = ?
GROUP BY c.id, c.name
`

	formattedStart := formatQueryDate(startDate, false)
	formattedEnd := formatQueryDate(endDate, true)

	err := config.DB.Raw(sqlQuery,
		companyID, formattedStart, formattedEnd, // channel_stats
		companyID, formattedStart, formattedEnd, // agent_stats
		companyID, formattedStart, formattedEnd, // oldest_open
		formattedStart, formattedEnd,            // cases LEFT JOIN
		companyID,                               // companies WHERE c.id = ?
	).Scan(&dashboard).Error

	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// By company and department
func (r *CaseDashboardRepository) GetDashboardByCompanyAndDepartment(companyID int64, departmentID int64, startDate, endDate string) (*models.CompanyDashboard, error) {
	if startDate == "" && endDate == "" {
		if err := r.refreshMaterializedWithDepartmentView(); err != nil {
			return nil, err
		}

		var dashboard models.CompanyDashboard
		if err := config.DB.
			Table("public.vw_case_dashboard_by_company_with_department").
			Where("company_id = ? AND department_id = ?", companyID, departmentID).
			First(&dashboard).Error; err != nil {
			return nil, err
		}
		return &dashboard, nil
	}

	// Dynamic query with date filter
	var dashboard models.CompanyDashboard
	sqlQuery := `
WITH channel_stats AS (
    SELECT cs_1.company_id,
        cs_1.department_id,
        ch.id AS channel_id,
        ch.name AS channel_name,
        count(cs_1.id) FILTER (WHERE (cs_1.status = 'open'::text)) AS open_cases,
        count(cs_1.id) FILTER (WHERE (cs_1.status = 'closed'::text)) AS closed_cases
    FROM public.cases cs_1
    LEFT JOIN public.channels ch ON ((ch.id = cs_1.channel_id))
    WHERE cs_1.company_id = ? AND cs_1.department_id = ? AND cs_1.created_at >= ? AND cs_1.created_at <= ?
    GROUP BY cs_1.company_id, cs_1.department_id, ch.id, ch.name
), agent_stats AS (
    SELECT cs_1.company_id,
        cs_1.department_id,
        u.id AS agent_id,
        u.full_name AS agent_name,
        count(cs_1.id) FILTER (WHERE (cs_1.status = 'open'::text)) AS open_cases,
        count(cs_1.id) FILTER (WHERE (cs_1.status = 'closed'::text)) AS closed_cases,
        round(avg((EXTRACT(epoch FROM (cs_1.closed_at - CASE WHEN cs_1.started_at IS NULL OR cs_1.started_at <= '1970-01-02'::timestamp THEN cs_1.created_at ELSE cs_1.started_at END)) / (3600)::numeric)), 2) AS avg_close_hours
    FROM public.cases cs_1
    LEFT JOIN public.users u ON ((u.id = cs_1.agent_id))
    WHERE (cs_1.agent_id IS NOT NULL) AND cs_1.company_id = ? AND cs_1.department_id = ? AND cs_1.created_at >= ? AND cs_1.created_at <= ?
    GROUP BY cs_1.company_id, cs_1.department_id, u.id, u.full_name
), oldest_open AS (
    SELECT cs_1.company_id,
        cs_1.department_id,
        cs_1.id AS case_id,
        cl.full_name AS client_name,
        cl.phone AS client_phone,
        cs_1.created_at,
        (SELECT m.created_at FROM public.messages m WHERE m.case_id = cs_1.id ORDER BY m.id DESC LIMIT 1) AS last_message_at,
        (SELECT m.sender_type FROM public.messages m WHERE m.case_id = cs_1.id ORDER BY m.id DESC LIMIT 1) AS last_message_sender_type
    FROM public.cases cs_1
    LEFT JOIN public.clients cl ON ((cl.id = cs_1.client_id))
    WHERE (cs_1.status = 'open'::text) AND cs_1.company_id = ? AND cs_1.department_id = ? AND cs_1.created_at >= ? AND cs_1.created_at <= ?
)
SELECT c.id AS company_id,
    c.name AS company_name,
    cs.department_id,
    count(cs.id) AS total_cases,
    count(cs.id) FILTER (WHERE (cs.status = 'open'::text)) AS open_cases,
    count(cs.id) FILTER (WHERE (cs.status = 'closed'::text)) AS closed_cases,
    count(cs.id) FILTER (WHERE ((cs.status = 'closed'::text) AND (date(cs.closed_at) = CURRENT_DATE))) AS closed_today,
    count(cs.id) FILTER (WHERE ((cs.status = 'open'::text) AND (date(cs.created_at) = CURRENT_DATE))) AS opened_today,
    count(cs.id) FILTER (WHERE (cs.status = 'cancelled'::text)) AS cancelled_cases,
    count(cs.id) FILTER (WHERE (
        cs.status IN ('open', 'in_progress')
        AND (
            SELECT m.sender_type
            FROM public.messages m
            WHERE m.case_id = cs.id
            ORDER BY m.id DESC
            LIMIT 1
        ) = 'client'
    )) AS unanswered_cases,
    count(cs.id) FILTER (WHERE (cs.agent_id IS NULL)) AS unassigned_agents,
    count(cs.id) FILTER (WHERE (cs.client_id IS NULL)) AS unassigned_clients,
    round(avg((EXTRACT(epoch FROM (cs.closed_at - CASE WHEN cs.started_at IS NULL OR cs.started_at <= '1970-01-02'::timestamp THEN cs.created_at ELSE cs.started_at END)) / (3600)::numeric)), 2) AS avg_close_hours,
    ( SELECT json_agg(json_build_object('channel_id', ch.channel_id, 'channel_name', ch.channel_name, 'open_cases', ch.open_cases, 'closed_cases', ch.closed_cases))
      FROM channel_stats ch
      WHERE (ch.company_id = c.id AND ch.department_id = cs.department_id)) AS cases_by_channel,
    ( SELECT json_agg(json_build_object('agent_id', a.agent_id, 'agent_name', a.agent_name, 'open_cases', a.open_cases, 'closed_cases', a.closed_cases, 'avg_close_hours', a.avg_close_hours))
      FROM agent_stats a
      WHERE (a.company_id = c.id AND a.department_id = cs.department_id)) AS cases_by_agent,
    ( SELECT json_agg(json_build_object('case_id', o.case_id, 'client_name', o.client_name, 'client_phone', o.client_phone, 'created_at', o.created_at, 'last_message_at', o.last_message_at, 'last_message_sender_type', o.last_message_sender_type) ORDER BY o.created_at)
      FROM oldest_open o
      WHERE (o.company_id = c.id AND o.department_id = cs.department_id)
      LIMIT 20) AS oldest_open_cases
FROM public.companies c
LEFT JOIN public.cases cs ON ((cs.company_id = c.id) AND cs.department_id = ? AND cs.created_at >= ? AND cs.created_at <= ?)
WHERE c.id = ?
GROUP BY c.id, c.name, cs.department_id
`

	formattedStart := formatQueryDate(startDate, false)
	formattedEnd := formatQueryDate(endDate, true)

	err := config.DB.Raw(sqlQuery,
		companyID, departmentID, formattedStart, formattedEnd, // channel_stats
		companyID, departmentID, formattedStart, formattedEnd, // agent_stats
		companyID, departmentID, formattedStart, formattedEnd, // oldest_open
		departmentID, formattedStart, formattedEnd,             // cases LEFT JOIN
		companyID,                                              // companies WHERE c.id = ?
	).Scan(&dashboard).Error

	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

func formatQueryDate(dateStr string, isEnd bool) string {
	if len(dateStr) == 10 { // YYYY-MM-DD
		if isEnd {
			return dateStr + " 23:59:59"
		}
		return dateStr + " 00:00:00"
	}
	return dateStr
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

// GetCasesByStatus obtiene los casos detallados filtrados por compañía, departamento, estado, búsqueda y fechas
func (r *CaseDashboardRepository) GetCasesByStatus(companyID int64, departmentID *int64, status string, search string, page int, limit int, startDate, endDate string) ([]models.CaseWithChannel, int64, error) {
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

	if startDate != "" {
		formattedStart := formatQueryDate(startDate, false)
		db = db.Where("created_at >= ?", formattedStart)
	}
	if endDate != "" {
		formattedEnd := formatQueryDate(endDate, true)
		db = db.Where("created_at <= ?", formattedEnd)
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
