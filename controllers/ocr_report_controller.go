package controllers

import (
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OcrReportController struct {
	repo *repository.OcrReportRepository
}

func NewOcrReportController() *OcrReportController {
	return &OcrReportController{
		repo: repository.NewOcrReportRepository(),
	}
}

func (c *OcrReportController) GetOcrReport(ctx *gin.Context) {
	companyIDParam := ctx.Param("company_id")
	companyID, err := strconv.ParseInt(companyIDParam, 10, 64)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	report, err := c.repo.GetOcrReport(companyID, startDate, endDate)
	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error al obtener el reporte de OCR", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Reporte de OCR obtenido correctamente", report, nil)
}
