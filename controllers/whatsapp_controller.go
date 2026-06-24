package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"harmony_api/config"
	"harmony_api/dto"
	"harmony_api/mapper"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/services"
	"harmony_api/utils"
	"harmony_api/ws"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type WhatsAppWebhookController struct {
	MessageEntry *MessageEntry
}

func NewWhatsAppWebhookController(hub *ws.Hub, ras *services.ReceiptAnalysisService) *WhatsAppWebhookController {
	return &WhatsAppWebhookController{
		MessageEntry: NewMessageEntry(hub, ras),
	}
}

func (w *WhatsAppWebhookController) Verify(c *gin.Context) {
	verifyToken := os.Getenv("WHATSAPP_VERIFY_TOKEN")
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	fmt.Println("=========================================")
	fmt.Println("🔍 WEBHOOK VERIFICATION (WHATSAPP):")
	fmt.Printf("   Mode: %s\n", mode)
	fmt.Printf("   Token: %s\n", token)
	fmt.Printf("   Challenge: %s\n", challenge)
	fmt.Printf("   Expected Token: %s\n", verifyToken)
	fmt.Println("=========================================")

	if mode == "subscribe" && token == verifyToken {
		fmt.Println("✅ Verificación exitosa")
		c.String(http.StatusOK, challenge)
	} else {
		fmt.Println("❌ Error de verificación: Token inválido")
		c.JSON(http.StatusForbidden, gin.H{"error": "Token inválido"})
	}
}

// -----------------------------------------------
//
//	ENDPOINT PARA REPROCESAR JSON DESDE N8N
//
// -----------------------------------------------
// func (w *WhatsAppWebhookController) ReceiveManual(c *gin.Context) {
// 	var raw []map[string]interface{} // 👈 porque n8n envía un ARRAY

// 	if err := c.ShouldBindJSON(&raw); err != nil {
// 		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido (array n8n)", nil, err)
// 		return
// 	}

// 	if len(raw) == 0 {
// 		utils.Respond(c, http.StatusBadRequest, false, "JSON vacío (n8n)", nil, nil)
// 		return
// 	}

// 	// n8n siempre manda el mensaje dentro de raw[0].body
// 	body, ok := raw[0]["body"]
// 	if !ok {
// 		utils.Respond(c, http.StatusBadRequest, false, "El objeto no contiene 'body'", nil, nil)
// 		return
// 	}

// 	// Convertir el body a JSON para rebindearlo
// 	bodyBytes, err := json.Marshal(body)
// 	if err != nil {
// 		utils.Respond(c, http.StatusInternalServerError, false, "Error procesando body", nil, err)
// 		return
// 	}

// 	// Bind al wrapper oficial
// 	var wrapper dto.N8nWhatsAppWrapper
// 	if err := json.Unmarshal(bodyBytes, &wrapper); err != nil {
// 		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido (wrapper n8n)", nil, err)
// 		return
// 	}

// 	// Convertimos wrapper → WhatsAppWebhookRequest normal
// 	req := dto.WhatsAppWebhookRequest{
// 		Entry: wrapper.Entry,
// 	}

// 	// Pasamos al mapper
// 	incoming := mapper.ParseWhatsAppToIncoming(req)

// 	repository :=
// 		repository.MessageRepository{}

// 	exists, err := repository.GetMessageControl(incoming.MessageID)
// 	if err != nil {
// 		utils.Respond(c, http.StatusInternalServerError, false, "Error verificando mensaje duplicado", nil, err)
// 		return
// 	}

// 	if exists {
// 		utils.Respond(c, http.StatusOK, true, "Mensaje duplicado ignorado", incoming, nil)
// 		return
// 	}

// 	// Procesar según tipo
// 	//w.dispatchIncomingMessage(c, incoming)
// 	// Procesar según tipo
// 	switch incoming.MessageType {

// 	case "text":
// 		newMessage, err := w.MessageEntry.processor.ProcessIncomingMessage(*incoming)
// 		if err != nil {
// 			utils.Respond(c, http.StatusInternalServerError, false, "Error reprocesando mensaje", nil, err)
// 			return
// 		}

// 		utils.Respond(c, http.StatusOK, true, "Mensaje reprocesado", newMessage, nil)

// 	case "image":
// 		w.MessageEntry.ReceiveImageMessageWebhookMedia(c)

// 	case "audio":
// 		w.MessageEntry.ReceiveAudioMessageWebhookMedia(c)

// 	default:
// 		utils.Respond(c, http.StatusOK, true, "Tipo de mensaje no soportado aún", incoming, nil)
// 	}
// }

func (w *WhatsAppWebhookController) ReceiveManual(c *gin.Context) {
	var raw []map[string]interface{} // n8n envía ARRAY

	if err := c.ShouldBindJSON(&raw); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido (array n8n)", nil, err)
		return
	}

	if len(raw) == 0 {
		utils.Respond(c, http.StatusBadRequest, false, "JSON vacío (n8n)", nil, nil)
		return
	}

	body, ok := raw[0]["body"]
	if !ok {
		utils.Respond(c, http.StatusBadRequest, false, "El objeto no contiene 'body'", nil, nil)
		return
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error procesando body", nil, err)
		return
	}

	var wrapper dto.N8nWhatsAppWrapper
	if err := json.Unmarshal(bodyBytes, &wrapper); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido (wrapper n8n)", nil, err)
		return
	}

	req := dto.WhatsAppWebhookRequest{
		Entry: wrapper.Entry,
	}

	incoming := mapper.ParseWhatsAppToIncoming(req)

	repo := repository.MessageRepository{}
	exists, err := repo.GetMessageControl(incoming.MessageID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error verificando mensaje duplicado", nil, err)
		return
	}

	if exists {
		utils.Respond(c, http.StatusOK, true, "Mensaje duplicado ignorado", incoming, nil)
		return
	}

	// 🔥 ÚNICO punto de entrada
	newMessage, err := w.MessageEntry.processor.ProcessIncomingMessage(*incoming)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error reprocesando mensaje", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Mensaje reprocesado correctamente", newMessage, nil)

	// 🔄 Refresh MV async
	mv := services.NewCaseMVRefreshService()
	mv.RefreshOnEvent("whatsapp_message_manual")
}

// POST: Recibir mensaje desde Meta
func (w *WhatsAppWebhookController) Receive(c *gin.Context) {
	var req dto.WhatsAppWebhookRequest
	raw, _ := c.GetRawData()

	// Reinyectar body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

	// LOG RAW JSON
	fmt.Println("=========================================")
	fmt.Println("📩 RAW WEBHOOK FROM META:")
	fmt.Println(string(raw))
	fmt.Println("=========================================")

	// Parse JSON Meta → DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido desde Meta", nil, err)
		return
	}

	// ---------------------------------------------------------
	// CAPTURA DE ESTADOS (SENT, DELIVERED, READ)
	// ---------------------------------------------------------

	fmt.Println("=========================================")
	fmt.Println("📩 RAW WEBHOOK FROM META:")
	fmt.Println(req.Entry)
	fmt.Println("=========================================")

	if len(req.Entry) > 0 && len(req.Entry[0].Changes) > 0 {
		change := req.Entry[0].Changes[0]

		// 1. CAPTURA DE CAMBIOS DE ESTADO DE PLANTILLAS EN META (APPROVED, REJECTED, etc.)
		if change.Field == "message_template_status_update" {
			event := change.Value.Event
			templateName := change.Value.MessageTemplateName
			templateIDStr := ""
			if change.Value.MessageTemplateID != nil {
				templateIDStr = fmt.Sprintf("%v", change.Value.MessageTemplateID)
			}

			statusMap := map[string]string{
				"APPROVED":         "approved",
				"PENDING_APPROVAL": "pending",
				"REJECTED":         "rejected",
				"PAUSED":           "paused",
				"DISABLED":         "disabled",
			}
			appStatus := "pending"
			if s, exists := statusMap[event]; exists {
				appStatus = s
			}

			updates := map[string]interface{}{
				"approval_status":  appStatus,
				"meta_template_id": templateIDStr,
			}
			if event == "REJECTED" && change.Value.Reason != "" {
				updates["rejection_reason"] = change.Value.Reason
			}

			// Localizar plantilla por nombre y actualizar
			var dbTemplate models.MessageTemplate
			if err := config.DB.Where("template_name = ?", templateName).First(&dbTemplate).Error; err == nil {
				config.DB.Model(&dbTemplate).Updates(updates)
				fmt.Printf("✅ Plantilla '%s' actualizada por Webhook a estado '%s'\n", templateName, appStatus)
			}

			utils.Respond(c, http.StatusOK, true, "Cambio de estado de plantilla procesado", nil, nil)
			return
		}

		// 2. CAPTURA DE ESTADOS DE MENSAJES (SENT, DELIVERED, READ)
		changeVal := change.Value
		if len(changeVal.Statuses) > 0 {
			repo := repository.MessageRepository{}

			for _, status := range changeVal.Statuses {
				fmt.Println("🔵 ESTADO WHATSAPP DETECTADO:")
				fmt.Printf("   ID: %s\n", status.ID)
				fmt.Printf("   Recipient: %s\n", status.RecipientID)
				fmt.Printf("   Status: %s\n", status.Status)
				fmt.Println("-----------------------------------------")

				// ------------------------------------------------------------------
				// LOGGING DE ESTADOS (NUEVO REQUERIMIENTO)
				// ------------------------------------------------------------------
				msgStatusLog := models.MessageStatus{
					ChannelMessageID: status.ID,
					MessageStatus:    status.Status,
					Applied:          false, // Por defecto false, o true si el update de abajo funciona
				}

				if err := config.DB.Create(&msgStatusLog).Error; err != nil {
					fmt.Printf("❌ Error guardando log de estado mensaje %s: %v\n", status.ID, err)
				} else {
					fmt.Printf("✅ Log de estado guardado ID: %d\n", msgStatusLog.ID)
				}

				// Actualizar estado en la BD
				if err := repo.UpdateMessageStatusByChannelID(status.ID, status.Status); err != nil {
					fmt.Printf("❌ Error actualizando estado mensaje %s: %v\n", status.ID, err)
				} else {
					fmt.Printf("✅ Estado mensaje %s actualizado a %s\n", status.ID, status.Status)

					// Marcar como aplicado
					// config.DB.Model(&msgStatusLog).Update("applied", true)
				}
			}
			// Si es solo status, no procesamos como mensaje entrante (text, image, etc)
			// Retornamos 200 OK para que Meta deje de reintentar
			utils.Respond(c, http.StatusOK, true, "Status recibido", nil, nil)
			return
		}
	}

	// Convertir DTO Meta → IncomingMessage unificado
	incoming := mapper.ParseWhatsAppToIncoming(req)

	repository :=
		repository.MessageRepository{}

	exists, err := repository.GetMessageControl(incoming.MessageID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error verificando mensaje duplicado", nil, err)
		return
	}

	if exists {
		utils.Respond(c, http.StatusOK, true, "Mensaje duplicado ignorado", incoming, nil)
		return
	}

	// AHORA aplicamos la misma lógica del MessageEntry, sin repetir código
	// switch incoming.MessageType {

	// case "text":
	// 	//w.forwardTextMessage(c, incoming)
	// 	newMessage, err := w.MessageEntry.processor.ProcessIncomingMessage(*incoming)
	// 	if err != nil {
	// 		utils.Respond(c, http.StatusInternalServerError, false, "Error procesando mensaje", nil, err)
	// 		return
	// 	}

	// 	utils.Respond(c, http.StatusOK, true, "Mensaje recibido", newMessage, nil)

	// case "image":
	// 	w.forwardImageMessage(c, incoming)

	// case "audio":
	// 	w.forwardAudioMessage(c, incoming)

	// default:
	// 	utils.Respond(c, http.StatusOK, true, "Tipo de mensaje no soportado aún", incoming, nil)
	// }

	switch incoming.MessageType {

	case "text", "image", "audio", "file":
		newMessage, err := w.MessageEntry.processor.ProcessIncomingMessage(*incoming)
		if err != nil {
			utils.Respond(c, http.StatusInternalServerError, false, "Error procesando mensaje", nil, err)
			return
		}

		utils.Respond(c, http.StatusOK, true, "Mensaje recibido", newMessage, nil)

	default:
		utils.Respond(c, http.StatusOK, true, "Tipo de mensaje no soportado aún", incoming, nil)
	}

	// 🔄 Refresh MV async
	mv := services.NewCaseMVRefreshService()
	mv.RefreshOnEvent("whatsapp_message")
}

// func (w *WhatsAppWebhookController) forwardTextMessage(c *gin.Context, input *models.IncomingMessage) {
// 	// Serializar para reutilizar el controller existente
// 	jsonBody, _ := json.Marshal(input)
// 	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

// 	w.MessageEntry.ReceiveMessageWebhook(c)
// }

// func (w *WhatsAppWebhookController) forwardImageMessage(c *gin.Context, input *models.IncomingMessage) {
// 	jsonBody, _ := json.Marshal(input)
// 	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

// 	w.MessageEntry.ReceiveImageMessageWebhookMedia(c)
// }

// func (w *WhatsAppWebhookController) forwardAudioMessage(c *gin.Context, input *models.IncomingMessage) {
// 	jsonBody, _ := json.Marshal(input)
// 	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

// 	w.MessageEntry.ReceiveAudioMessageWebhookMedia(c)
// }

func (w *WhatsAppWebhookController) dispatchIncomingMessage(c *gin.Context, incoming *models.IncomingMessage) {
	jsonBody, _ := json.Marshal(incoming)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

	switch incoming.MessageType {
	case "text":
		w.MessageEntry.ReceiveMessageWebhook(c)

	case "image":
		w.MessageEntry.ReceiveImageMessageWebhookMedia(c)

	case "audio":
		w.MessageEntry.ReceiveAudioMessageWebhookMedia(c)

	default:
		utils.Respond(c, http.StatusOK, true, "Tipo de mensaje no soportado aún", incoming, nil)
	}
}
