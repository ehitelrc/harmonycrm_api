package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterChannelRoutes(r *gin.RouterGroup) {
	controller := &controllers.ChannelController{}

	r.GET("/channels", controller.GetAll)
	r.GET("/channels/:id", controller.GetByID)
	r.POST("/channels", controller.Create)
	r.PUT("/channels", controller.Update)
	r.DELETE("/channels/:id", controller.Delete)

	// Channel Whatsapp template routes
	r.POST("/channels/whatsapp/templates", controller.CreateWhatsappTemplate)
	r.GET("/channels/whatsapp/templates/company/:company_id", controller.GetWhatsappTemplatesByCompanyID)

	// Get channel integration by company_id

	r.GET("/channels/integrations/whatsapp/company/:company_id", controller.GetChannelWhatsappIntegrationsByCompanyID)

	// Get templates by channel_integration_id
	r.GET("/channels/whatsapp/templates/integration/:channel_integration_id", controller.GetWhatsappTemplatesByChannelIntegrationID)

}
