package controllers

import (
	"harmony_api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UnreconciledPaymentsController struct {
	repo *repository.UnreconciledPaymentsRepository
}

func NewUnreconciledPaymentsController() *UnreconciledPaymentsController {
	return &UnreconciledPaymentsController{
		repo: repository.NewUnreconciledPaymentsRepository(),
	}
}

func (c *UnreconciledPaymentsController) GetUnreconciledPayments(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	payments, err := c.repo.GetUnreconciledPayments(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al buscar pagos no conciliados",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payments,
	})
}
