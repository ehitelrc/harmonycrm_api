package controllers

import (
	"net/http"

	"harmony_api/services"
	"harmony_api/utils"

	"github.com/gin-gonic/gin"
)

type OCRController struct {
	service *services.OCRService
}

func NewOCRController(service *services.OCRService) *OCRController {
	return &OCRController{service: service}
}

type OCRRequest struct {
	Base64 string `json:"base64" binding:"required"`
}

func (ctrl *OCRController) OCR(c *gin.Context) {
	var req OCRRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Debe enviar un base64 válido", nil, err)
		return
	}

	text, err := ctrl.service.ProcessBase64(req.Base64)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error realizando OCR", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "OCR realizado con éxito", map[string]string{
		"text": text,
	}, nil)
}
