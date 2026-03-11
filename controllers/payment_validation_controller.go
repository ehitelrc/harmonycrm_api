package controllers

import (
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentValidationController struct {
}

func NewPaymentValidationController() *PaymentValidationController {
	return &PaymentValidationController{}
}

// GetPaymentValidations fetches the payment validation view data based on query parameters.
func (c *PaymentValidationController) GetPaymentValidations(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	caseIdStr := ctx.Query("case_id")
	contractIdStr := ctx.Query("contract_id")

	var caseId, contractId int
	if caseIdStr != "" {
		parsed, err := strconv.Atoi(caseIdStr)
		if err == nil {
			caseId = parsed
		}
	}
	if contractIdStr != "" {
		parsed, err := strconv.Atoi(contractIdStr)
		if err == nil {
			contractId = parsed
		}
	}

	repo := repository.NewPaymentValidationRepository()
	data, err := repo.GetPaymentValidations(startDate, endDate, caseId, contractId)

	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error getting payment validations", nil, err)
		return
	}

	utils.Respond(ctx, http.StatusOK, true, "Payment validations retrieved successfully", data, nil)
}

// GetReceipt fetches the receipt base64 string directly from the DB.
func (c *PaymentValidationController) GetReceipt(ctx *gin.Context) {
	erpIDStr := ctx.Param("erp_id")
	if erpIDStr == "" {
		utils.Respond(ctx, http.StatusBadRequest, false, "erp_id is required", nil, nil)
		return
	}

	erpID, err := strconv.ParseInt(erpIDStr, 10, 64)
	if err != nil {
		utils.Respond(ctx, http.StatusBadRequest, false, "invalid erp_id", nil, nil)
		return
	}

	repo := repository.NewPaymentValidationRepository()
	base64Str, err := repo.GetReceiptBase64(erpID)

	if err != nil {
		utils.Respond(ctx, http.StatusInternalServerError, false, "Error getting receipt base64", nil, err)
		return
	}

	// Because image base64 formats can vary and we might just want to return the raw string or wrapped JSON
	utils.Respond(ctx, http.StatusOK, true, "Receipt fetched successfully", map[string]string{
		"receipt_base64": base64Str,
	}, nil)
}
