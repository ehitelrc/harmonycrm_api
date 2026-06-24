package controllers

import (
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CaseDashboardController struct {
	repo *repository.CaseDashboardRepository
}

func NewCaseDashboardController() *CaseDashboardController {
	return &CaseDashboardController{repo: repository.NewCaseDashboardRepository()}
}

func (c *CaseDashboardController) GetCompanyDashboard(ctx *gin.Context) {
	companyIDParam := ctx.Param("company_id")

	companyID, err := strconv.Atoi(companyIDParam)

	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el ID de la empresa", nil, err)
		return
	}

	dashboard, err := c.repo.GetDashboardByCompany(int64(companyID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el dashboard de la empresa", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Dashboard de la empresa obtenido correctamente", dashboard, nil)

}

func (c *CaseDashboardController) GetDepartmentDashboard(ctx *gin.Context) {
	companyIDParam := ctx.Param("company_id")
	departmentIDParam := ctx.Param("department_id")

	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el ID de la empresa", nil, err)
		return
	}

	departmentID, err := strconv.Atoi(departmentIDParam)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el ID del departamento", nil, err)
		return
	}

	dashboard, err := c.repo.GetDashboardByCompanyAndDepartment(int64(companyID), int64(departmentID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el dashboard del departamento", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Dashboard del departamento obtenido correctamente", dashboard, nil)

}

func (c *CaseDashboardController) GetUserDashboard(ctx *gin.Context) {
	companyIDParam := ctx.Param("company_id")
	userIDParam := ctx.Param("user_id")

	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el ID de la empresa", nil, err)
		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el ID del usuario", nil, err)
		return
	}

	dashboard, err := c.repo.GetDashboardByCompanyAndUser(int64(companyID), int64(userID))
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el dashboard del usuario", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Dashboard del usuario obtenido correctamente", dashboard, nil)

}

func (c *CaseDashboardController) GetCasesByStatus(ctx *gin.Context) {
	companyIDParam := ctx.Param("company_id")

	companyID, err := strconv.ParseInt(companyIDParam, 10, 64)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "ID de compañía inválido", nil, err)
		return
	}

	status := ctx.Query("status")
	search := ctx.Query("search")
	deptIDStr := ctx.Query("department_id")
	var departmentID *int64
	if deptIDStr != "" {
		if dId, err := strconv.ParseInt(deptIDStr, 10, 64); err == nil && dId > 0 {
			departmentID = &dId
		}
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	cases, total, err := c.repo.GetCasesByStatus(companyID, departmentID, status, search, page, limit)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener casos", nil, err)
		return
	}

	response := gin.H{
		"cases": cases,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	utils.Respond(ctx, http.StatusOK, true, "Casos obtenidos correctamente", response, nil)
}
