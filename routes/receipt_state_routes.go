package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func ReceiptStateRoutes(router *gin.RouterGroup, ctrl *controllers.ReceiptStateController) {
	r := router.Group("/case-receipts")
	{
		r.GET("/new", ctrl.ListNew)
		r.PUT("/:id/read", ctrl.MarkAsRead)
		r.PUT("/:id/processed", ctrl.MarkAsProcessed)
	}
}
