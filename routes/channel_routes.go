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

}
