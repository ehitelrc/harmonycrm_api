package routes

import (
	"harmony_api/controllers"
	"harmony_api/ws"

	"github.com/gin-gonic/gin"
)

func RegisterWhatsappWebHookRoutes(r *gin.RouterGroup, hub *ws.Hub) {

	webhook := controllers.NewWhatsAppWebhookController(hub)

	r.GET("/webhook/whatsapp", webhook.Verify)
	r.POST("/webhook/whatsapp", webhook.Receive)

	// Endpoint manual (n8n / pruebas / Postman)
	r.POST("/webhook/whatsapp/manual", webhook.ReceiveManual)
}
