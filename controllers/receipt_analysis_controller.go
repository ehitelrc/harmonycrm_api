package controllers

import (
	"net/http"
	"time"

	"harmony_api/repository"
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

type AnalyzeAndSaveReceiptRequest struct {
	Base64    string `json:"base64" binding:"required"`
	CaseID    uint   `json:"case_id" binding:"required"`
	CreatedAt string `json:"created_at"` // optional format YYYY-MM-DD or full timestamp
}

func (ctrl *ReceiptAnalysisController) AnalyzeAndSaveReceipt(c *gin.Context) {
	var req AnalyzeAndSaveReceiptRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Debe enviar 'base64' y 'case_id' válidos", nil, err)
		return
	}

	result, err := ctrl.service.AnalyzeFromBase64(c.Request.Context(), req.Base64, &req.CaseID, true)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error analizando recibo", nil, err)
		return
	}

	if result == nil {
		utils.Respond(c, http.StatusBadRequest, false, "La imagen procesada no es un recibo", nil, nil)
		return
	}

	// Como el método subyacente no guarda a pesar del flag save=true (por cómo está implementado internamente),
	// instanciamos el repositorio y lo guardamos tal como se hace en el flujo de los mensajes.
	receiptRepo := repository.NewReceiptRepository()

	var customCreatedAt *time.Time
	if req.CreatedAt != "" {
		// Attempt to parse YYYY-MM-DD
		parsedTime, errParse := time.Parse("2006-01-02", req.CreatedAt)
		if errParse == nil {
			customCreatedAt = &parsedTime
		} else {
			// Fallback attempt for RFC3339 if needed
			if parsedTime, errParse := time.Parse(time.RFC3339, req.CreatedAt); errParse == nil {
				customCreatedAt = &parsedTime
			}
		}
	}

	record, err := receiptRepo.SaveReceiptResult(result, req.CaseID, customCreatedAt)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error guardando recibo en base de datos", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Recibo analizado y guardado con éxito", map[string]interface{}{
		"record": record,
		"extracted_data": result,
	}, nil)
}
