package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func ReceiptAnalysisRoutes(router *gin.RouterGroup, ctrl *controllers.ReceiptAnalysisController) {
	r := router.Group("/receipts")
	{
		// POST /api/receipts/analyze
		r.POST("/analyze", ctrl.AnalyzeReceipt)

		// POST /api/receipts/analyze-and-save
		r.POST("/analyze-and-save", ctrl.AnalyzeAndSaveReceipt)
	}
}
