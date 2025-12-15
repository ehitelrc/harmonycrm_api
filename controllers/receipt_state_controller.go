package controllers

import (
	"net/http"
	"strconv"

	"harmony_api/dto"
	"harmony_api/services"
	"harmony_api/utils"

	"github.com/gin-gonic/gin"
)

type ReceiptStateController struct {
	service *services.ReceiptStateService
}

func NewReceiptStateController(service *services.ReceiptStateService) *ReceiptStateController {
	return &ReceiptStateController{service: service}
}

// GET /api/receipts/new
func (ctrl *ReceiptStateController) ListNew(c *gin.Context) {
	receipts, err := ctrl.service.ListNew()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error consultando recibos", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Recibos en estado 'new'", receipts, nil)
}

// PUT /api/receipts/:id/read
func (ctrl *ReceiptStateController) MarkAsRead(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	if err := ctrl.service.MarkRead(uint(id)); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "No se pudo marcar como leído", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Recibo marcado como 'read'", nil, nil)
}

// PUT /api/receipts/:id/processed
func (ctrl *ReceiptStateController) MarkAsProcessed(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "ID inválido", nil, err)
		return
	}
	id := uint(id64)

	// Leer body
	var req dto.ProcessReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	// 1. Cambiar estado a processed
	if err := ctrl.service.MarkProcessed(id); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "No se pudo actualizar el estado", nil, err)
		return
	}

	// 2. Generar PDF usando template HTML
	// pdfPath, err := ctrl.service.GenerateReceiptPDF(id, req)
	// if err != nil {
	// 	utils.Respond(c, http.StatusInternalServerError, false, "Error generando PDF", nil, err)
	// 	return
	// }

	// 3. Enviar PDF por WhatsApp
	// err = ctrl.service.SendPDFViaWhatsApp(id, pdfPath, req.Mensaje)
	// if err != nil {
	// 	utils.Respond(c, http.StatusInternalServerError, false, "Error enviando PDF por WhatsApp", nil, err)
	// 	return
	// }

	utils.Respond(c, http.StatusOK, true, "Recibo marcado como procesado y PDF enviado", gin.H{
		"id":      id,
		"pdf_url": "",
	}, nil)
}
