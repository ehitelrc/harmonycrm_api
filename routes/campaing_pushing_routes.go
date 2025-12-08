package routes

import (
	"harmony_api/controllers"
	"harmony_api/ws"

	"github.com/gin-gonic/gin"
)

func RegisterCampaignPushingRoutes(r *gin.RouterGroup, hub *ws.Hub) {
	ctrl := controllers.NewCampaignPushingController(hub)

	// Register campaign pushing
	r.POST("/campaigns/whatsapp/push/register", ctrl.RegisterWhatsappCampaignPush)

	// Send template message
	r.POST("/campaigns/whatsapp/send-template/template/:template_id/case/:case_id", ctrl.SendWhatsappTemplateMessage)

	// New Case from template
	r.POST("/campaigns/whatsapp/new-case/template", ctrl.CreateNewWhatsappCaseFromTemplate)

}
