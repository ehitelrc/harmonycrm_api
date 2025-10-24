package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"harmony_api/ws"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MessageEntry struct {
	hub *ws.Hub
}

func NewMessageEntry(hub *ws.Hub) *MessageEntry {
	return &MessageEntry{hub: hub}
}

type WSMessage struct {
	Type   string      `json:"type"` // "new_message"
	CaseID uint        `json:"case_id"`
	Data   interface{} `json:"data"` // el mensaje recién guardado o un DTO
}

func (m *MessageEntry) ReceiveMessageWebhook(c *gin.Context) {
	var input models.IncomingMessage

	// Leer el cuerpo sin procesar
	rawData, _ := c.GetRawData()
	fmt.Println("Raw JSON recibido:", string(rawData))

	// Reinyectar el cuerpo para poder hacer el binding después de leerlo
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	repository :=
		repository.MessageRepository{}

	newMessage, err := repository.CreateMessage(input)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al procesar el mensaje", nil, err)
		return
	}

	// Broadcast WS (si tenemos case_id)
	if newMessage.CaseID != 0 && m.hub != nil {
		payload, _ := json.Marshal(WSMessage{
			Type:   "new_message",
			CaseID: uint(newMessage.CaseID),
			Data:   newMessage, // o arma un DTO si prefieres
		})
		channel := "case:" + strconv.Itoa(int(newMessage.CaseID))
		m.hub.BroadcastJSON(channel, payload)
	}

	utils.Respond(c, http.StatusOK, true, "Mensaje recibido correctamente", input, nil)
}

// Get cases without agent assigned by company_id

func (m *MessageEntry) GetCasesWithoutAgentByCompanyID(c *gin.Context) {
	companyID := c.Param("company_id")

	companyIDInt, err := strconv.Atoi(companyID)

	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	repository := repository.MessageRepository{}

	cases, err := repository.GetUnassignedCasesByCompanyID(companyIDInt)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los casos sin agente asignado", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Casos sin agente asignado obtenidos correctamente!", cases, nil)
}

func (m *MessageEntry) ReceiveImageMessageWebhookMedia(c *gin.Context) {
	var input models.IncomingMessage

	// Leer el cuerpo sin procesar
	rawData, _ := c.GetRawData()
	fmt.Println("Raw JSON recibido:", string(rawData))

	// Reinyectar el cuerpo para poder hacer el binding después de leerlo
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	// Get channel integration
	channelRepository := repository.ChannelRepository{}

	channnel, err := channelRepository.GetChannelIntegrationByAppIdentifier(input.RecipientID)

	if err != nil || channnel == nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener la integración del canal", nil, err)
		return
	}

	wm_utils := utils.WSMediaMessage{}

	mediaUrl := fmt.Sprintf("https://graph.facebook.com/v23.0/%s", input.MediaID)

	_, resourceData, error := wm_utils.GetMediaData(mediaUrl, *channnel)

	if error != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los datos del medio", nil, error)
		return
	}

	completeData := "data:" + input.MIMEType + ";base64," + resourceData

	input.Base64Content = completeData

	repository :=
		repository.MessageRepository{}

	newMessage, err := repository.CreateMessage(input)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al procesar el mensaje", nil, err)
		return
	}

	// Broadcast WS (si tenemos case_id)
	if newMessage.CaseID != 0 && m.hub != nil {
		payload, _ := json.Marshal(WSMessage{
			Type:   "new_message",
			CaseID: uint(newMessage.CaseID),
			Data:   input, // o arma un DTO si prefieres
		})
		channel := "case:" + strconv.Itoa(int(newMessage.CaseID))
		m.hub.BroadcastJSON(channel, payload)
	}

	utils.Respond(c, http.StatusOK, true, "Mensaje recibido correctamente", input, nil)
}

func (m *MessageEntry) ReceiveAudioMessageWebhookMedia(c *gin.Context) {
	var input models.IncomingMessage

	// Leer el cuerpo sin procesar
	rawData, _ := c.GetRawData()
	fmt.Println("Raw JSON recibido:", string(rawData))

	// Reinyectar el cuerpo para poder hacer el binding después de leerlo
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	// Get channel integration
	channelRepository := repository.ChannelRepository{}

	channnel, err := channelRepository.GetChannelIntegrationByAppIdentifier(input.RecipientID)

	if err != nil || channnel == nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener la integración del canal", nil, err)
		return
	}

	wm_utils := utils.WSMediaMessage{}

	mediaUrl := fmt.Sprintf("https://graph.facebook.com/v23.0/%s", input.MediaID)

	_, resourceData, error := wm_utils.GetMediaData(mediaUrl, *channnel)

	if error != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los datos del medio", nil, error)
		return
	}

	completeData := "data:" + input.MIMEType + ";base64," + resourceData

	input.Base64Content = completeData

	repository :=
		repository.MessageRepository{}

	newMessage, err := repository.CreateMessage(input)

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al procesar el mensaje", nil, err)
		return
	}

	// Broadcast WS (si tenemos case_id)
	if newMessage.CaseID != 0 && m.hub != nil {
		payload, _ := json.Marshal(WSMessage{
			Type:   "new_message",
			CaseID: uint(newMessage.CaseID),
			Data:   input, // o arma un DTO si prefieres
		})
		channel := "case:" + strconv.Itoa(int(newMessage.CaseID))
		m.hub.BroadcastJSON(channel, payload)
	}

	utils.Respond(c, http.StatusOK, true, "Mensaje recibido correctamente", input, nil)
}

func (m *MessageEntry) GetActiveCasesByAgentID(c *gin.Context) {
	agentID := c.Param("agent_id")

	repository := repository.MessageRepository{}

	activeCases, err := repository.GetActiveCasesByAgentID(agentID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los casos activos", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Casos activos obtenidos correctamente!", activeCases, nil)
}

func (m *MessageEntry) GetMessagesByCaseID(c *gin.Context) {
	caseID := c.Param("case_id")

	repository := repository.MessageRepository{}

	messages, err := repository.GetMessagesByCaseID(caseID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los mensajes", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Mensajes obtenidos correctamente!", messages, nil)
}

func (m *MessageEntry) SendMessageToPlatform(c *gin.Context) {
	var input models.AgentMessage

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	channelRepository := repository.ChannelRepository{}

	channelIntegration, err := channelRepository.GetChannerlByCaseID(input.CaseID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener la integración del canal", nil, err)
		return
	}

	if channelIntegration != nil {
		if channelIntegration.ChannelCode == "messenger" && input.MessageType == "text" {
			err := m.DispatchTextMessage(channelIntegration, input)

			if err != nil {
				utils.Respond(c, http.StatusInternalServerError, false, "Error al enviar el mensaje", nil, err)
				return
			}
		} else if channelIntegration.ChannelCode == "whatsapp" && input.MessageType == "text" {
			err := m.DispatchWhatsappTextMessage(channelIntegration, input)
			if err != nil {
				utils.Respond(c, http.StatusInternalServerError, false, "Error al enviar el mensaje", nil, err)
				return
			}
		} else if channelIntegration.ChannelCode == "whatsapp" && input.MessageType == "image" {

			media_id, err := uploadBase64Image(*channelIntegration.AppIdentifier, *channelIntegration.AccessToken, input.Base64Content)

			if err != nil {
				utils.Respond(c, http.StatusInternalServerError, false, "Error al subir la imagen", nil, err)
				return
			}

			//Extract mime type
			parts := strings.SplitN(input.Base64Content, ";", 2)
			//remove "data:"
			mimeType := parts[0]
			//remove "base64,"
			mimeType = strings.TrimPrefix(mimeType, "data:")
			mimeType = strings.TrimPrefix(mimeType, "base64,")

			recipientId := channelIntegration.SenderID

			input.MIMEType = mimeType

			err = m.sendWhatsAppImage(*channelIntegration.AppIdentifier, *channelIntegration.AccessToken, *recipientId, media_id, input.TextMessage)
			if err != nil {
				utils.Respond(c, http.StatusInternalServerError, false, "Error al enviar el mensaje", nil, err)
				return
			}
		}
	}

	repository := repository.MessageRepository{}

	if err := repository.SendMessageToPlatform(input); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al enviar el mensaje", nil, err)
		return
	}

	// Consumir

	payload, _ := json.Marshal(WSMessage{
		Type:   "new_message",
		CaseID: uint(input.CaseID),
		Data:   input, // o arma un DTO si prefieres
	})
	channel := "case:" + strconv.Itoa(int(input.CaseID))
	m.hub.BroadcastJSON(channel, payload)

	utils.Respond(c, http.StatusOK, true, "Mensaje enviado correctamente", input, nil)
}

func (m *MessageEntry) AssignCaseToClient(c *gin.Context) {
	var input models.AssignCaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	msgRepository := repository.MessageRepository{}

	if err := msgRepository.AssignCaseToClient(input); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al asignar el caso", nil, err)
		return
	}

	// Get case
	caseData, err := msgRepository.GetCaseByID(input.CaseID)

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener el caso", nil, err)
		return
	}

	channelId, _ := strconv.ParseUint(caseData.ChannelID, 10, 64)

	// Client social account
	clientSocialAccountRepo := repository.ClientSocialAccountRepository{}

	clientSocialAccount, err := clientSocialAccountRepo.GetByChannelAndExternal(uint(channelId), caseData.SenderId)

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener la cuenta social del cliente", nil, err)
		return
	}

	// If not exists then create
	if clientSocialAccount == nil {

		newClientSocialAccount := models.ClientSocialAccount{
			ClientID:   input.ClientID,
			ChannelID:  uint(channelId),
			ExternalID: caseData.SenderId,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := clientSocialAccountRepo.Create(&newClientSocialAccount)

		if err != nil {
			utils.Respond(c, http.StatusInternalServerError, false, "Error al crear la cuenta social del cliente", nil, err)
			return
		}
	}

	utils.Respond(c, http.StatusOK, true, "Caso asignado correctamente", input, nil)

}

func (m *MessageEntry) AddCaseNote(c *gin.Context) {
	var input models.CaseNote

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Datos inválidos", nil, err)
		return
	}

	repository := repository.MessageRepository{}

	if err := repository.AddCaseNote(input); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al agregar la nota del caso", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Nota del caso agregada correctamente", input, nil)
}

func (m *MessageEntry) GetCaseNotesByCaseID(c *gin.Context) {
	caseID := c.Param("case_id")

	repository := repository.MessageRepository{}
	notes, err := repository.GetCaseNotesByCaseID(caseID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener las notas del caso", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Notas del caso obtenidas correctamente", notes, nil)
}

func (m *MessageEntry) DispatchTextMessage(channelIntegration *models.VWCaseChannelIntegration, message models.AgentMessage) error {
	url := channelIntegration.WebhookURL
	accessToken := channelIntegration.AccessToken
	me := channelIntegration.AppIdentifier
	recipientId := channelIntegration.SenderID

	// Construir el payload según el formato del canal
	payload := map[string]string{
		"access_token": strOrEmpty(accessToken),
		"me":           strOrEmpty(me),
		"recipient_id": strOrEmpty(recipientId),
		"message_text": message.TextMessage,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializando payload: %w", err)
	}

	// Crear request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}

	// Headers
	req.Header.Set("Content-Type", "application/json")
	// if accessToken != nil && *accessToken != "" {
	// 	req.Header.Set("Authorization", "Bearer "+*accessToken)
	// }

	// Cliente HTTP con timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("error en respuesta webhook: %s", resp.Status)
	}

	return nil
}

func (m *MessageEntry) DispatchWhatsappTextMessage(channelIntegration *models.VWCaseChannelIntegration, message models.AgentMessage) error {
	url := channelIntegration.WebhookURL
	//url := "https://ehitelrc.app.n8n.cloud/webhook-test/6b2b114c-f863-44b6-8ab6-a80968c24d82"
	accessToken := channelIntegration.AccessToken
	me := channelIntegration.AppIdentifier
	recipientId := channelIntegration.SenderID

	// Construir el payload según el formato del canal
	payload := map[string]string{
		"access_token": strOrEmpty(accessToken),
		"me":           strOrEmpty(me),
		"to":           strOrEmpty(recipientId),
		"recipient_id": strOrEmpty(recipientId),
		"message_text": message.TextMessage,
		"message_type": "text",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializando payload: %w", err)
	}

	// Crear request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}

	// Headers
	req.Header.Set("Content-Type", "application/json")
	// if accessToken != nil && *accessToken != "" {
	// 	req.Header.Set("Authorization", "Bearer "+*accessToken)
	// }

	// Cliente HTTP con timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("error en respuesta webhook: %s", resp.Status)
	}

	return nil
}

func uploadBase64Image(phoneNumberID, accessToken, base64DataParam string) (string, error) {
	// 1️⃣ Limpieza de prefijo
	parts := strings.SplitN(base64DataParam, ",", 2)
	base64Data := base64DataParam
	if len(parts) == 2 {
		base64Data = parts[1]
	}

	// 2️⃣ Sanitización
	re := regexp.MustCompile(`[^A-Za-z0-9+/=]`)
	base64Data = re.ReplaceAllString(base64Data, "")
	base64Data = strings.TrimSpace(base64Data)

	// 3️⃣ Detección de tipo MIME
	mimeType := "image/jpeg"
	ext := "jpg"
	if strings.Contains(base64DataParam, "image/png") {
		mimeType = "image/png"
		ext = "png"
	}

	// 4️⃣ Decodificación Base64
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("❌ error decodificando base64: %w", err)
	}

	// 5️⃣ Guardar archivo temporal
	os.MkdirAll("./tmp", 0755)
	tempPath := fmt.Sprintf("./tmp/debug_image_upload.%s", ext)
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return "", fmt.Errorf("❌ error guardando temporal: %w", err)
	}
	fmt.Println("🧩 Archivo temporal:", tempPath)

	// 6️⃣ Crear multipart con header correcto
	file, err := os.Open(tempPath)
	if err != nil {
		return "", fmt.Errorf("❌ error abriendo archivo temporal: %w", err)
	}
	defer file.Close()

	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/media", phoneNumberID)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// ✅ Parte con MIMEHeader explícito (clave del error original)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(tempPath)))
	h.Set("Content-Type", mimeType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return "", fmt.Errorf("❌ error creando part con header: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return "", fmt.Errorf("❌ error copiando bytes: %w", err)
	}

	writer.WriteField("messaging_product", "whatsapp")
	writer.Close()

	// 7️⃣ Enviar request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("❌ error creando request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("❌ error enviando request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("📡 Respuesta API:", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("❌ error API (%d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error parseando JSON: %w", err)
	}

	id, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("no se encontró el campo 'id' en la respuesta: %s", string(respBody))
	}

	fmt.Println("📦 ID del archivo subido:", id)
	return id, nil
}

func (m *MessageEntry) sendWhatsAppImage(phoneNumberID, accessToken, to, mediaID, caption string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/messages", phoneNumberID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "image",
		"image": map[string]interface{}{
			"id":      mediaID,
			"caption": caption,
		},
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("❌ error creando request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("❌ error enviando request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("📡 Respuesta de envío:", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("❌ error API (%d): %s", resp.StatusCode, string(respBody))
	}

	fmt.Println("✅ Mensaje enviado correctamente.")
	return nil
}

// controllers/message_entry.go
func (m *MessageEntry) AssignCaseToCampaign(c *gin.Context) {
	var req models.AssignCaseToCampaignRequest

	// Bind del JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	// Si changed_by no viene en el body, lo puedes obtener del contexto/auth
	if req.ChangedBy == 0 {
		req.ChangedBy = c.GetInt("user_id") // ejemplo si lo guardas en middleware
	}

	repo := repository.MessageRepository{}
	if err := repo.AssignCaseToCampaign(req.CaseID, req.CampaignID, req.ChangedBy); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al asignar el caso a la campaña", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Caso asignado a la campaña correctamente", nil, nil)
}

func (m *MessageEntry) AssignCaseToDepartment(c *gin.Context) {
	var req models.AssignCaseToDepartmentRequest

	// Bind del JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	// Si changed_by no viene en el body, lo puedes obtener del contexto/auth
	if req.ChangedBy == 0 {
		req.ChangedBy = c.GetInt("user_id") // ejemplo si lo guardas en middleware
	}

	repo := repository.MessageRepository{}
	if err := repo.AssignCaseToDepartment(req.CaseID, req.DepartmentID, req.ChangedBy); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al asignar el caso al departamento", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Caso asignado al departamento correctamente", nil, nil)
}

func (m *MessageEntry) AssignCaseToAgent(c *gin.Context) {
	var req models.AssignCaseToAgentRequest

	// Bind del JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	// Si changed_by no viene en el body, lo puedes obtener del contexto/auth
	if req.ChangedBy == 0 {
		req.ChangedBy = c.GetInt("user_id") // ejemplo si lo guardas en middleware
	}

	repo := repository.MessageRepository{}
	if err := repo.AssignCaseToAgent(req.CaseID, req.AgentID, req.ChangedBy, req.DepartmentID); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al asignar el caso al agente", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Caso asignado al agente correctamente", nil, nil)
}

// GetCurrentCaseFunnel
func (m *MessageEntry) GetCurrentCaseFunnel(c *gin.Context) {
	// Obtener case_id desde los parámetros de la URL
	caseIDStr := c.Param("case_id")
	caseID, err := strconv.Atoi(caseIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "case_id inválido", nil, err)
		return
	}

	repo := repository.MessageRepository{}

	funnel, err := repo.GetCurrentCaseFunnel(caseID)
	if err != nil {
		utils.Respond(c, http.StatusOK, false, "Error al obtener el funnel del caso", nil, err)
		return
	}

	if funnel.CaseID == 0 {

		currentCase, err := repo.GetCaseByID(uint(caseID))

		if err != nil {
			utils.Respond(c, http.StatusOK, false, "Error al obtener el caso", nil, err)
			return
		}

		if currentCase == nil {
			utils.Respond(c, http.StatusOK, false, "Caso no encontrado", nil, fmt.Errorf("caso no encontrado"))
			return
		}

		// Get funnel

		funnelRepo := repository.FunnelRepository{}

		currentFunnel, err := funnelRepo.GetByID(uint(currentCase.FunnelID))

		if err != nil {
			utils.Respond(c, http.StatusOK, false, "Error al obtener el funnel", nil, err)
			return
		}

		funnelID := int(currentCase.FunnelID)

		funnel.CaseID = int(currentCase.ID)
		funnel.FunnelID = &funnelID
		funnel.FunnelName = &currentFunnel.Name

	}

	utils.Respond(c, http.StatusOK, true, "Funnel del caso obtenido correctamente", funnel, nil)
}

func (m *MessageEntry) SetCaseFunnelStage(c *gin.Context) {
	var req models.MoveCaseStagePayload

	// Bind del JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	caseFunnelStage := models.CaseFunnel{

		CaseID:      req.CaseID,
		FunnelID:    &req.FunnelID,
		FromStageID: req.FromStageID,
		ToStageID:   &req.ToStageID,
		Note:        req.Note,
		ChangedBy:   *req.ChangedBy, // Aquí deberías obtener el ID del usuario que hace el cambio
		ChangedAt:   time.Now(),
		Action:      "move",
	}

	repo := repository.MessageRepository{}
	if err := repo.SetCaseFunnelStage(caseFunnelStage); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al actualizar la etapa del funnel del caso", nil, err)
		return

	}

	// Si changed_by no viene en el body, lo puedes obtener del contexto/auth

	utils.Respond(c, http.StatusOK, true, "Etapa del funnel del caso actualizada correctamente", nil, nil)
}

func (m *MessageEntry) CloseCase(c *gin.Context) {
	var req models.CaseCloseRequest

	// Bind del JSON

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	repo := repository.MessageRepository{}
	if err := repo.CloseCase(req); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al cerrar el caso", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Caso cerrado correctamente", nil, nil)
}

func (m *MessageEntry) GetCaseGeneralInformation(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	campaignIDStr := c.Param("campaign_id")
	campaignID, err := strconv.Atoi(campaignIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "campaign_id inválido", nil, err)
		return
	}

	stageIDStr := c.Param("stage_id")
	stageID, err := strconv.Atoi(stageIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "stage_id inválido", nil, err)
		return
	}

	repo := repository.MessageRepository{}
	cases, err := repo.GetCaseGeneralInformation(uint(companyID), uint(campaignID), uint(stageID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener la información general de los casos", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Información general de los casos obtenida correctamente", cases, nil)

}

func (m *MessageEntry) GetCaseGeneralInformationByCompanyCampaignAgent(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	campaignIDStr := c.Param("campaign_id")
	campaignID, err := strconv.Atoi(campaignIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "campaign_id inválido", nil, err)
		return
	}

	agentIDStr := c.Param("agent_id")
	agentID, err := strconv.Atoi(agentIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "agent_id inválido", nil, err)
		return
	}

	channelIntegrationIDStr := c.Param("channel_integration_id")
	channelIntegrationID, err := strconv.Atoi(channelIntegrationIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "channel_integration_id inválido", nil, err)
		return
	}

	repo := repository.MessageRepository{}
	cases, err := repo.GetCaseGeneralInformationByCompanyCampaignAgent(uint(companyID), uint(campaignID), uint(agentID), uint(channelIntegrationID))
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener la información general de los casos", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Información general de los casos obtenida correctamente", cases, nil)

}

// Helper para manejar punteros string nulos
func strOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
