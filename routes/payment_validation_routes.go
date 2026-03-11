package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentValidationRoutes(router *gin.RouterGroup) {
	controller := controllers.NewPaymentValidationController()

	group := router.Group("/payment-validations")
	{
		group.GET("", controller.GetPaymentValidations)
		group.GET("/receipt/:erp_id", controller.GetReceipt)
	}
}
