package controllers

import (
	"bytes"
	"encoding/json"
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

	if mode == "subscribe" && token == verifyToken {
		c.String(http.StatusOK, challenge)
	} else {
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

	// Parse JSON Meta → DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido desde Meta", nil, err)
		return
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
