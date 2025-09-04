package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterClientSocialAccountRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewClientSocialAccountController()

	// Básicos
	r.GET("/client-social-accounts", ctrl.GetAll)
	r.GET("/client-social-accounts/:id", ctrl.GetByID)
	r.POST("/client-social-accounts", ctrl.Create)
	r.PUT("/client-social-accounts", ctrl.Update) // objeto completo con id
	r.DELETE("/client-social-accounts/:id", ctrl.Delete)

	// Filtros
	r.GET("/client-social-accounts/client/:client_id", ctrl.GetByClient)
	r.GET("/client-social-accounts/channel/:channel_id", ctrl.GetByChannel)
	r.GET("/client-social-accounts/channel/:channel_id/external/:external_id", ctrl.GetByChannelAndExternal)

	// Activación / desactivación rápida
	r.PATCH("/client-social-accounts/:id/activate", ctrl.Activate)
	r.PATCH("/client-social-accounts/:id/deactivate", ctrl.Deactivate)
}
