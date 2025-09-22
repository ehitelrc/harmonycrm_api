package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/utils"
	"harmony_api/ws"
	"io"
	"net/http"
	"strconv"
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

	wm_utils := utils.WSMediaMessage{}

	mediaUrl := fmt.Sprintf("https://graph.facebook.com/v23.0/%s", input.MediaID)

	_, resourceData, error := wm_utils.GetMediaData(mediaUrl)

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

	wm_utils := utils.WSMediaMessage{}

	mediaUrl := fmt.Sprintf("https://graph.facebook.com/v23.0/%s", input.MediaID)

	_, resourceData, error := wm_utils.GetMediaData(mediaUrl)

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

	repository := repository.MessageRepository{}

	if err := repository.AssignCaseToClient(input); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al asignar el caso", nil, err)
		return
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
	if err := repo.AssignCaseToAgent(req.CaseID, req.AgentID, req.ChangedBy); err != nil {
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
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener el funnel del caso", nil, err)
		return
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
		ID:          uint(req.ToStageID),
		CaseID:      req.CaseID,
		FunnelID:    req.FunnelID,
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

// Helper para manejar punteros string nulos
func strOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
