package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUnreconciledPaymentsRoutes(router *gin.RouterGroup) {
	controller := controllers.NewUnreconciledPaymentsController()

	group := router.Group("/unreconciled-payments")
	{
		group.GET("", controller.GetUnreconciledPayments)
	}
}
