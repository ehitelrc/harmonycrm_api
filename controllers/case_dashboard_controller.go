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
