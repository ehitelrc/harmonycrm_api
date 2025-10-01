package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCampaignPushingRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewCampaignPushingController()

	// Register campaign pushing
	r.POST("/campaigns/whatsapp/push/register", ctrl.RegisterWhatsappCampaignPush)

}
