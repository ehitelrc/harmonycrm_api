package controllers

import (
	"fmt"
	"harmony_api/dto"
	"harmony_api/mapper"
	"harmony_api/services"
	"harmony_api/ws"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type MessengerWebhookController struct {
	hub *ws.Hub
	ras *services.ReceiptAnalysisService
}

func NewMessengerWebhookController(
	hub *ws.Hub,
	ras *services.ReceiptAnalysisService,
) *MessengerWebhookController {
	return &MessengerWebhookController{
		hub: hub,
		ras: ras,
	}
}

func (c *MessengerWebhookController) Verify(ctx *gin.Context) {
	mode := ctx.Query("hub.mode")
	token := ctx.Query("hub.verify_token")
	challenge := ctx.Query("hub.challenge")
	verifyToken := os.Getenv("MESSENGER_VERIFY_TOKEN")

	fmt.Println("=========================================")
	fmt.Println("🔍 WEBHOOK VERIFICATION (MESSENGER):")
	fmt.Printf("   Mode: %s\n", mode)
	fmt.Printf("   Token: %s\n", token)
	fmt.Printf("   Challenge: %s\n", challenge)
	fmt.Printf("   Expected Token: %s\n", verifyToken)
	fmt.Println("=========================================")

	if mode == "subscribe" && token == verifyToken {
		fmt.Println("✅ Verificación exitosa (Messenger)")
		ctx.String(http.StatusOK, challenge)
		return
	}

	fmt.Println("❌ Error de verificación (Messenger): Token inválido")
	ctx.AbortWithStatus(http.StatusForbidden)
}

func (c *MessengerWebhookController) Receive(ctx *gin.Context) {
	var payload dto.MessengerWebhookRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	incomingMessages := mapper.ParseMessengerWebhook(payload)

	processor := services.NewMessageProcessor(
		c.hub,
		c.ras,
	)

	for _, msg := range incomingMessages {
		_, _ = processor.ProcessIncomingMessage(msg)
	}

	ctx.Status(http.StatusOK)
}
