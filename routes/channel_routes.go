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
	r.PUT("/channels/whatsapp/templates", controller.UpdateWhatsappTemplate)
	r.DELETE("/channels/whatsapp/templates/:template_id", controller.DeleteWhatsappTemplate)

	// Get channel integration by company_id

	r.GET("/channels/integrations/whatsapp/company/:company_id", controller.GetChannelWhatsappIntegrationsByCompanyID)

	// Add integration to channel
	r.POST("/channels/integrations", controller.AddIntegrationToChannel)

	//Update channel integration
	r.PUT("/channels/integrations", controller.UpdateChannelIntegration)

	// Integrations by company and channel_id
	r.GET("/channels/integrations/company/:company_id/channel/:channel_id", controller.GetChannelIntegrationsByCompanyAndChannelID)

	// Get templates by channel_integration_id
	r.GET("/channels/whatsapp/templates/integration/:channel_integration_id", controller.GetWhatsappTemplatesByChannelIntegrationID)

}
