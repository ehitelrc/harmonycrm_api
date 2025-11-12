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
