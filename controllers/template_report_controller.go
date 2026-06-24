package controllers

import (
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateReportController struct {
	repo *repository.TemplateReportRepository
}

func NewTemplateReportController() *TemplateReportController {
	return &TemplateReportController{
		repo: repository.NewTemplateReportRepository(),
	}
}

func (c *TemplateReportController) GetTemplateReport(ctx *gin.Context) {
	companyIDParam := ctx.Param("company_id")
	companyID, err := strconv.ParseInt(companyIDParam, 10, 64)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	report, err := c.repo.GetTemplateReport(companyID)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el reporte de plantillas", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Reporte de plantillas obtenido correctamente", report, nil)
}
