package routes

import (
	"harmony_api/controllers"
	"harmony_api/services"
	"harmony_api/ws"

	"github.com/gin-gonic/gin"
)

func RegisterMessengerWebHookRoutes(
	r *gin.RouterGroup,
	hub *ws.Hub,
	ras *services.ReceiptAnalysisService,
) {

	webhook := controllers.NewMessengerWebhookController(hub, ras)

	r.GET("/webhook/messenger", webhook.Verify)
	r.POST("/webhook/messenger", webhook.Receive)

}
