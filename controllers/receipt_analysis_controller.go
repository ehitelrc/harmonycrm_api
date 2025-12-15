package controllers

import (
	"net/http"

	"harmony_api/services"
	"harmony_api/utils"

	"github.com/gin-gonic/gin"
)

type ReceiptAnalysisController struct {
	service *services.ReceiptAnalysisService
}

func NewReceiptAnalysisController(service *services.ReceiptAnalysisService) *ReceiptAnalysisController {
	return &ReceiptAnalysisController{service: service}
}

type AnalyzeReceiptRequest struct {
	Base64 string `json:"base64" binding:"required"`
}

func (ctrl *ReceiptAnalysisController) AnalyzeReceipt(c *gin.Context) {
	var req AnalyzeReceiptRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Debe enviar el campo 'base64'", nil, err)
		return
	}

	// Llamamos al servicio SOLO para analizar, sin guardar
	result, err := ctrl.service.AnalyzeFromBase64(c.Request.Context(), req.Base64, nil, false)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error analizando recibo", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Recibo analizado con éxito", result, nil)
}
