package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"harmony_api/dto"
	"harmony_api/models"
	"harmony_api/repository"
	"harmony_api/services"
	"harmony_api/utils"
	"harmony_api/ws"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// type MessageEntry struct {
// 	hub                    *ws.Hub
// 	receiptAnalysisService *services.ReceiptAnalysisService
// }

type MessageEntry struct {
	hub                    *ws.Hub
	processor              *services.MessageProcessor
	receiptAnalysisService *services.ReceiptAnalysisService
}

// func NewMessageEntry(hub *ws.Hub, ras *services.ReceiptAnalysisService) *MessageEntry {
// 	return &MessageEntry{
// 		hub:                    hub,
// 		receiptAnalysisService: ras,
// 	}
// }

func NewMessageEntry(
	hub *ws.Hub,
	ras *services.ReceiptAnalysisService,
) *MessageEntry {

	return &MessageEntry{
		hub:                    hub,
		processor:              services.NewMessageProcessor(hub, ras),
		receiptAnalysisService: ras,
	}
}

type WSMessage struct {
	Type   string      `json:"type"` // "new_message"
	CaseID uint        `json:"case_id"`
	Data   interface{} `json:"data"` // el mensaje recién guardado o un DTO
}

// func (m *MessageEntry) ReceiveMessageWebhook(c *gin.Context) {
// 	var input models.IncomingMessage

// 	// Leer el cuerpo sin procesar
// 	rawData, _ := c.GetRawData()
// 	fmt.Println("Raw JSON recibido:", string(rawData))

// 	// Reinyectar el cuerpo para poder hacer el binding después de leerlo
// 	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
// 		return
// 	}

// 	repository :=
// 		repository.MessageRepository{}

// 	newMessage, err := repository.CreateMessage(input)
// 	if err != nil {
// 		utils.Respond(c, http.StatusInternalServerError, false, "Error al procesar el mensaje", nil, err)
// 		return
// 	}

// 	// Broadcast WS (si tenemos case_id)
// 	if newMessage.CaseID != 0 && m.hub != nil {
// 		payload, _ := json.Marshal(WSMessage{
// 			Type:   "new_message",
// 			CaseID: uint(newMessage.CaseID),
// 			Data:   newMessage, // o arma un DTO si prefieres
// 		})
// 		channel := "case:" + strconv.Itoa(int(newMessage.CaseID))
// 		m.hub.BroadcastJSON(channel, payload)
// 		if newMessage.AgentID != nil {
// 			m.hub.BroadcastJSON("agent:"+strconv.Itoa(int(*newMessage.AgentID)), payload)
// 		}
// 	}

// 	fmt.Println("Nuevo mensaje guardado con ID:", newMessage.ID)

// 	utils.Respond(c, http.StatusOK, true, "Mensaje recibido correctamente", input, nil)

// }

func (m *MessageEntry) ReceiveMessageWebhook(c *gin.Context) {
	var input models.IncomingMessage

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	newMessage, err := m.processor.ProcessIncomingMessage(input)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error procesando mensaje", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Mensaje recibido", newMessage, nil)

	// 🔄 Refresh MV async
	mv := services.NewCaseMVRefreshService()
	mv.RefreshOnEvent("whatsapp_message")
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

// By company and department unassigned cases
//api.GET("/entry/unassigned_cases/company/:company_id/department/:department_id", controller.GetCasesWithoutAgentByCompanyAndDepartmentID)

func (m *MessageEntry) GetCasesWithoutAgentByCompanyAndDepartmentID(c *gin.Context) {
	companyID := c.Param("company_id")
	departmentID := c.Param("department_id")
	companyIDInt, err := strconv.Atoi(companyID)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}
	departmentIDInt, err := strconv.Atoi(departmentID)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "department_id inválido", nil, err)
		return
	}

	repository := repository.MessageRepository{}
	cases, err := repository.GetUnassignedCasesByCompanyAndDepartmentID(companyIDInt, departmentIDInt)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los casos sin agente asignado", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Casos sin agente asignado obtenidos correctamente!", cases, nil)
}

// Get open cases by company_id and department_id
func (m *MessageEntry) GetOpenCasesByCompanyAndDepartmentID(c *gin.Context) {
	companyID := c.Param("company_id")
	departmentID := c.Param("department_id")
	companyIDInt, err := strconv.Atoi(companyID)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}
	departmentIDInt, err := strconv.Atoi(departmentID)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "department_id inválido", nil, err)
		return
	}

	repository := repository.MessageRepository{}
	openCases, err := repository.GetOpenCasesByCompanyAndDepartmentID(companyIDInt, departmentIDInt)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los casos abiertos", nil, err)
		return
	}
	utils.Respond(c, http.StatusOK, true, "Casos abiertos obtenidos correctamente!", openCases, nil)
}

func (m *MessageEntry) GetOpenCasesByCompanyAndDepartmentIDV2(c *gin.Context) {

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	departmentID, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "department_id inválido", nil, err)
		return
	}

	// query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if limit <= 0 {
		limit = 50
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if cached, found := utils.GetOpenCasesFromCache(
		uint(companyID),
		uint(departmentID),
		page,
		limit,
	); found {
		utils.Respond(c, http.StatusOK, true, "Casos abiertos (cache)", cached, nil)
		return
	}

	repo := repository.MessageRepository{}

	items, total, err := repo.GetOpenCasesByCompanyAndDepartmentIDPaged(
		companyID,
		departmentID,
		limit,
		offset,
	)

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los casos abiertos (v2)", nil, err)
		return
	}

	payload := gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	utils.CacheOpenCases(uint(companyID), uint(departmentID), page, limit, payload)

	utils.Respond(c, http.StatusOK, true, "Casos abiertos", payload, nil)

}

func (m *MessageEntry) GetOpenCasesStatsByCompanyAndDepartmentV2(c *gin.Context) {

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "company_id inválido", nil, err)
		return
	}

	departmentID, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "department_id inválido", nil, err)
		return
	}

	if cached, found := utils.GetOpenCasesStatsFromCache(
		uint(companyID),
		uint(departmentID),
	); found {
		utils.Respond(c, http.StatusOK, true, "Stats (cache)", cached, nil)
		return
	}

	repo := repository.MessageRepository{}
	total, assigned, unassigned, err :=
		repo.GetOpenCasesStatsByCompanyAndDepartment(companyID, departmentID)

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error obteniendo stats", nil, err)
		return
	}

	data := gin.H{
		"total_open": total,
		"assigned":   assigned,
		"unassigned": unassigned,
	}

	utils.CacheOpenCasesStats(uint(companyID), uint(departmentID), data)

	utils.Respond(c, http.StatusOK, true, "Stats", data, nil)

	// utils.Respond(c, http.StatusOK, true, "Stats de casos abiertos (v2)", gin.H{
	// 	"total_open": total,
	// 	"assigned":   assigned,
	// 	"unassigned": unassigned,
	// }, nil)
}

// Materialized
func (m *MessageEntry) GetOpenCasesMV(c *gin.Context) {
	companyID, _ := strconv.Atoi(c.Param("company_id"))
	departmentID, _ := strconv.Atoi(c.Param("department_id"))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	fmt.Printf(
		"[MV] company=%d dept=%d limit=%d offset=%d page=%d\n",
		companyID, departmentID, limit, offset, page,
	)

	repo := repository.MessageRepository{}
	cases, err := repo.GetOpenCasesMV(
		uint(companyID),
		uint(departmentID),
		limit,
		offset,
	)

	if err != nil {
		utils.Respond(c, 500, false, "Error obteniendo casos (MV)", nil, err)
		return
	}

	utils.Respond(c, 200, true, "Casos obtenidos (MV)", cases, nil)
}

func (m *MessageEntry) GetCaseStats(c *gin.Context) {
	companyID, _ := strconv.Atoi(c.Param("company_id"))
	departmentID, _ := strconv.Atoi(c.Param("department_id"))

	repo := repository.MessageRepository{}
	stats, err := repo.GetCaseStats(uint(companyID), uint(departmentID))
	if err != nil {
		utils.Respond(c, 500, false, "Error obteniendo stats", nil, err)
		return
	}

	utils.Respond(c, 200, true, "Stats obtenidas", stats, nil)
}

// MarkMessagesAsReadByCaseID
func (m *MessageEntry) MarkMessagesAsReadByCaseID(c *gin.Context) {
	caseID := c.Param("case_id")
	repository := repository.MessageRepository{}

	err := repository.MarkMessagesAsReadByCaseID(caseID)

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al marcar los mensajes como leídos", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Mensajes marcados como leídos correctamente", nil, nil)

}

// CreateNewCaseFromTemplate creates a new case from a template
func (m *MessageEntry) CreateNewCaseFromTemplate(c *gin.Context) {
	var requestBody dto.NewWhatsappCaseFromTemplateRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Invalid request body", nil, err)
		return
	}

	repo := repository.MessageRepository{}
	caseID, err := repo.NewCaseFromTemplate(requestBody, m.hub)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error creating new case from template", nil, err)
		return
	}

	if caseID != 0 && m.hub != nil {
		payload, _ := json.Marshal(map[string]interface{}{
			"type":    "new_message",
			"case_id": caseID,
			"data":    nil, // O un objeto vacío, la recarga del lado del FE se basa en el evento
		})

		channel := "case:" + strconv.Itoa(int(caseID))
		m.hub.BroadcastJSON(channel, payload)

		if requestBody.AgentID != 0 {
			m.hub.BroadcastJSON("agent:"+strconv.Itoa(requestBody.AgentID), payload)
		}
	}

	utils.Respond(c, http.StatusOK, true, "Nuevo caso creado correctamente", map[string]interface{}{
		"case_id": caseID,
	}, nil)
}

// SendTemplateToCase envía un template a un caso existente usando message_templates
func (m *MessageEntry) SendTemplateToCase(c *gin.Context) {
	templateID, err := strconv.Atoi(c.Param("template_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "template_id inválido", nil, err)
		return
	}

	caseID, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "case_id inválido", nil, err)
		return
	}

	repo := repository.MessageRepository{}
	newMessage, err := repo.SendTemplateToCase(templateID, caseID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error enviando template al caso", nil, err)
		return
	}

	// Broadcast WS
	if caseID != 0 && m.hub != nil {
		payload, _ := json.Marshal(map[string]interface{}{
			"type":    "new_message",
			"case_id": caseID,
			"data":    newMessage,
		})

		channel := "case:" + strconv.Itoa(caseID)
		m.hub.BroadcastJSON(channel, payload)

		if newMessage != nil && newMessage.AgentID != nil {
			m.hub.BroadcastJSON("agent:"+strconv.Itoa(int(*newMessage.AgentID)), payload)
		}
	}

	utils.Respond(c, http.StatusOK, true, "Template enviado correctamente", nil, nil)
}

// func (m *MessageEntry) ReceiveImageMessageWebhookMedia(c *gin.Context) {
// 	var input models.IncomingMessage

// 	// Leer el cuerpo sin procesar
// 	rawData, _ := c.GetRawData()
// 	fmt.Println("Raw JSON recibido:", string(rawData))

// 	// Reinyectar el cuerpo para poder hacer el binding después de leerlo
// 	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
// 		return
// 	}

// 	// Get channel integration
// 	channelRepository := repository.ChannelRepository{}

// 	channnel, err := channelRepository.GetChannelIntegrationByAppIdentifier(input.RecipientID)

// 	if err != nil || channnel == nil {
// 		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener la integración del canal", nil, err)
// 		return
// 	}

// 	wm_utils := utils.WSMediaMessage{}

// 	mediaUrl := fmt.Sprintf("https://graph.facebook.com/v23.0/%s", input.MediaID)

// 	_, resourceData, error := wm_utils.GetMediaData(mediaUrl, *channnel)

// 	if error != nil {
// 		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los datos del medio", nil, error)
// 		return
// 	}

// 	completeData := "data:" + input.MIMEType + ";base64," + resourceData

// 	input.Base64Content = completeData

// 	msgRepo :=
// 		repository.MessageRepository{}

// 	newMessage, err := msgRepo.CreateMessage(input)
// 	if err != nil {
// 		utils.Respond(c, http.StatusInternalServerError, false, "Error al procesar el mensaje", nil, err)
// 		return
// 	}

// 	// Broadcast WS (si tenemos case_id)
// 	if newMessage.CaseID != 0 && m.hub != nil {
// 		payload, _ := json.Marshal(WSMessage{
// 			Type:   "new_message",
// 			CaseID: uint(newMessage.CaseID),
// 			Data:   input, // o arma un DTO si prefieres
// 		})
// 		channel := "case:" + strconv.Itoa(int(newMessage.CaseID))
// 		m.hub.BroadcastJSON(channel, payload)
// 	}

// 	utils.Respond(c, http.StatusOK, true, "Mensaje recibido correctamente", input, nil)

// 	// Guardar información del recibo

// 	//---------------------------------------
// 	// ANALIZAR RECIBO CON OCR + OPENAI
// 	//---------------------------------------

// 	if !channnel.AnalyzeIncomingImages {
// 		fmt.Println("ℹ️ Análisis de imágenes entrantes deshabilitado para esta integración.")
// 		return
// 	}

// 	go func(input models.IncomingMessage, newMessage *models.Message) {

// 		// Base64 sin prefijo data:
// 		b64 := input.Base64Content
// 		if idx := strings.Index(b64, ","); idx != -1 {
// 			b64 = b64[idx+1:]
// 		}

// 		// Obtener caseID
// 		caseID := uint(newMessage.CaseID)

// 		// Ejecutar OCR + IA
// 		result, err := m.receiptAnalysisService.AnalyzeFromBase64(
// 			c,
// 			b64,
// 			&caseID,
// 			true, // es mensaje entrante
// 		)
// 		if err != nil {
// 			fmt.Println("❌ Error analizando recibo:", err)
// 			return
// 		}

// 		if result == nil {
// 			fmt.Println("ℹ️ La imagen no es un recibo.")
// 			return
// 		}

// 		// Guardar en la base de datos
// 		receiptRepo := repository.NewReceiptRepository()

// 		record, err := receiptRepo.SaveReceiptResult(result, caseID)
// 		if err != nil {
// 			fmt.Println("❌ Error guardando recibo:", err)
// 			return
// 		}

// 		fmt.Println("✅ Recibo guardado con ID:", record.ID)

// 	}(input, newMessage)
// }

func (m *MessageEntry) ReceiveImageMessageWebhookMedia(c *gin.Context) {
	var input models.IncomingMessage

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	newMessage, err := m.processor.ProcessIncomingMessage(input)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error procesando imagen", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Imagen recibida", newMessage, nil)
}

// func (m *MessageEntry) ReceiveAudioMessageWebhookMedia(c *gin.Context) {
// 	var input models.IncomingMessage

// 	// Leer el cuerpo sin procesar
// 	rawData, _ := c.GetRawData()
// 	fmt.Println("Raw JSON recibido:", string(rawData))

// 	// Reinyectar el cuerpo para poder hacer el binding después de leerlo
// 	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
// 		return
// 	}

// 	// Get channel integration
// 	channelRepository := repository.ChannelRepository{}

// 	channnel, err := channelRepository.GetChannelIntegrationByAppIdentifier(input.RecipientID)

// 	if err != nil || channnel == nil {
// 		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener la integración del canal", nil, err)
// 		return
// 	}

// 	wm_utils := utils.WSMediaMessage{}

// 	mediaUrl := fmt.Sprintf("https://graph.facebook.com/v23.0/%s", input.MediaID)

// 	_, resourceData, error := wm_utils.GetMediaData(mediaUrl, *channnel)

// 	if error != nil {
// 		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los datos del medio", nil, error)
// 		return
// 	}

// 	completeData := "data:" + input.MIMEType + ";base64," + resourceData

// 	input.Base64Content = completeData

// 	repository :=
// 		repository.MessageRepository{}

// 	newMessage, err := repository.CreateMessage(input)

// 	if err != nil {
// 		utils.Respond(c, http.StatusInternalServerError, false, "Error al procesar el mensaje", nil, err)
// 		return
// 	}

// 	// Broadcast WS (si tenemos case_id)
// 	if newMessage.CaseID != 0 && m.hub != nil {
// 		payload, _ := json.Marshal(WSMessage{
// 			Type:   "new_message",
// 			CaseID: uint(newMessage.CaseID),
// 			Data:   input, // o arma un DTO si prefieres
// 		})
// 		channel := "case:" + strconv.Itoa(int(newMessage.CaseID))
// 		m.hub.BroadcastJSON(channel, payload)
// 	}

// 	utils.Respond(c, http.StatusOK, true, "Mensaje recibido correctamente", input, nil)
// }

func (m *MessageEntry) ReceiveAudioMessageWebhookMedia(c *gin.Context) {
	var input models.IncomingMessage

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "JSON inválido", nil, err)
		return
	}

	newMessage, err := m.processor.ProcessIncomingMessage(input)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error procesando audio", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Audio recibido", newMessage, nil)
}

func (m *MessageEntry) GetActiveCasesByAgentID(c *gin.Context) {
	agentID := c.Param("agent_id")

	intAgentID, err := strconv.Atoi(agentID)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "agent_id inválido", nil, err)
		return
	}

	if cached, found := utils.GetActiveCasesByAgentFromCache(uint(intAgentID)); found {
		utils.Respond(c, http.StatusOK, true, "Casos activos (cache)", cached, nil)
		return
	}

	repository := repository.MessageRepository{}

	activeCases, err := repository.GetActiveCasesByAgentID(agentID)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los casos activos", nil, err)
		return
	}

	utils.CacheActiveCasesByAgent(uint(intAgentID), activeCases)

	utils.Respond(c, http.StatusOK, true, "Casos activos obtenidos correctamente!", activeCases, nil)
}

func (m *MessageEntry) GetMessagesByCaseID(c *gin.Context) {
	caseID := c.Param("case_id")

	intCaseID, err := strconv.Atoi(caseID)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "case_id inválido", nil, err)
		return
	}

	if cached, found := utils.GetFirstMessagesByCaseFromCache(uint(intCaseID)); found {
		utils.Respond(c, http.StatusOK, true, "Mensajes (cache)", cached, nil)
		fmt.Println("✅ Mensajes servidos desde cache para case_id:", caseID)
		return
	}

	repository := repository.MessageRepository{}

	messages, err := repository.GetMessagesByCaseID(caseID)

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al obtener los mensajes", nil, err)
		return
	}

	utils.CacheFirstMessagesByCase(uint(intCaseID), messages)

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

		repo := repository.MessageRepository{}

		open, err := repo.IsMessengerWindowOpen(input.CaseID)

		if err != nil {
			utils.Respond(c, http.StatusInternalServerError, false, "Error verificando ventana de mensajería", nil, err)
			return
		}

		if channelIntegration.ChannelCode == "messenger" && !open {

			// 1️⃣ Marcar error funcional
			input.HasError = true
			input.MessageError = "La ventana de mensajería de Messenger ha expirado (24h)"

			// 2️⃣ Guardar el mensaje igualmente (traza)
			repoMsg := repository.MessageRepository{}
			_ = repoMsg.SendMessageToPlatform(input)

			// 3️⃣ Responder al frontend (mensaje existe, pero falló)
			utils.Respond(
				c,
				http.StatusOK, // ⚠️ negocio OK, envío NO
				true,
				"Mensaje registrado pero no enviado (ventana de 24h expirada)",
				input,
				nil,
			)
			return
		}

		// if channelIntegration.ChannelCode == "messenger" && input.MessageType == "text" {
		// 	err := m.DispatchTextMessage(channelIntegration, input)

		// 	if err != nil {
		// 		utils.Respond(c, http.StatusInternalServerError, false, "Error al enviar el mensaje", nil, err)
		// 		return
		// 	}
		if channelIntegration.ChannelCode == "messenger" && input.MessageType == "text" {

			_, err := m.SendMessengerTextDirect(
				*channelIntegration.AccessToken,
				*channelIntegration.SenderID,
				input.TextMessage,
			)

			if err != nil {
				input.HasError = true
				input.MessageError = err.Error()
			} else {
				input.HasError = false
				input.MessageError = ""
			}
		} else if channelIntegration.ChannelCode == "messenger" && input.MessageType == "image" {

			imageURL, err := utils.SaveBase64ImageAndGetURL(
				input.Base64Content,
				os.Getenv("PUBLIC_BASE_URL"),
			)

			if err != nil {
				fmt.Println("❌ Error guardando imagen:", err)
				input.HasError = true
				input.MessageError = err.Error()
			} else {

				err = m.SendMessengerImage(
					*channelIntegration.AccessToken,
					*channelIntegration.SenderID,
					imageURL,
				)

				if err != nil {
					fmt.Println("❌ Error enviando imagen a Messenger:", err)
					input.HasError = true
					input.MessageError = err.Error()
				} else {
					input.HasError = false
					input.MessageError = ""
				}
			}
		} else if channelIntegration.ChannelCode == "messenger" && input.MessageType == "audio" {

			// 🎧 Si viene webm → convertir a ogg (Messenger NO reproduce webm)
			if input.MIMEType == "audio/webm" {

				convertedBase64, newMime, err := convertWebMToOgg(input.Base64Content)
				if err != nil {
					input.HasError = true
					input.MessageError = err.Error()
				} else {
					input.Base64Content = convertedBase64
					input.MIMEType = newMime // audio/ogg
					input.FileName = "audio.ogg"
				}
			}

			// Si hubo error en conversión, no seguimos
			if input.HasError {
				// no hacemos return, dejamos que se guarde con error
			} else {

				audioURL, err := utils.SaveBase64AudioAndGetURL(
					input.Base64Content,
					input.MIMEType,
					os.Getenv("PUBLIC_BASE_URL"),
				)

				if err != nil {
					input.HasError = true
					input.MessageError = err.Error()
				} else {

					err = m.SendMessengerAudio(
						*channelIntegration.AccessToken,
						*channelIntegration.SenderID,
						audioURL,
					)

					if err != nil {
						input.HasError = true
						input.MessageError = err.Error()
					} else {
						input.HasError = false
						input.MessageError = ""
					}
				}
			}
		} else if channelIntegration.ChannelCode == "whatsapp" && input.MessageType == "text" {
			//err := m.DispatchWhatsappTextMessage(channelIntegration, input)
			wamid, apiResponse, err := m.SendWhatsAppTextDirect(
				*channelIntegration.AppIdentifier,
				*channelIntegration.AccessToken,
				*channelIntegration.SenderID,
				input.TextMessage,
			)

			fmt.Println("✅ WAMID:", wamid)
			fmt.Println("✅ API Response:", apiResponse)

			if err != nil {

				// Generar log visible para debug
				fmt.Println("*********************************************************")
				fmt.Println("** ERROR ENVIANDO MENSAJE DE TEXTO DIRECTO:", err.Error())
				fmt.Println("*********************************************************")

				input.HasError = true
				input.MessageError = err.Error()

			} else {
				// Todo salió bien

				input.HasError = false
				input.MessageError = ""
				input.ChannelMessageID = wamid

			}

			// Notificar Frontend por WebSocket
			// payload, _ := json.Marshal(WSMessage{
			// 	Type:   "new_message",
			// 	CaseID: uint(input.CaseID),
			// 	Data:   input,
			// })
			// channel := "case:" + strconv.Itoa(int(input.CaseID))
			// m.hub.BroadcastJSON(channel, payload)

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

			wamid, err := m.sendWhatsAppImage(*channelIntegration.AppIdentifier, *channelIntegration.AccessToken, *recipientId, media_id, input.TextMessage)
			if err != nil {
				utils.Respond(c, http.StatusInternalServerError, false, "Error al enviar el mensaje", nil, err)
				return
			}

			input.HasError = false
			input.MessageError = ""
			input.ChannelMessageID = wamid

		} else if channelIntegration.ChannelCode == "whatsapp" && input.MessageType == "file" {

			mediaID, err := uploadBase64File(
				*channelIntegration.AppIdentifier,
				*channelIntegration.AccessToken,
				input.Base64Content, // ← RAW BASE64
				input.MIMEType,      // ← MIME EXACTO DEL FRONT
				input.FileName,      // ← NOMBRE DEL ARCHIVO
			)
			if err != nil {
				utils.Respond(c, http.StatusInternalServerError, false, "Error al subir archivo", nil, err)
				return
			}

			//sendWhatsAppDocument(appID, token, recipient, mediaID, caption string)

			id, err := m.sendWhatsAppDocument(
				*channelIntegration.AppIdentifier,
				*channelIntegration.AccessToken,
				*channelIntegration.SenderID,
				mediaID,
				input.TextMessage,
			)
			if err == nil {
				fmt.Println("✅ Documento enviado correctamente. ID:", id)
			}
			if err != nil {
				utils.Respond(c, http.StatusInternalServerError, false, "Error al enviar archivo", nil, err)
				return
			}

			input.HasError = false
			input.MessageError = ""
			input.ChannelMessageID = id

		} else if channelIntegration.ChannelCode == "whatsapp" && input.MessageType == "audio" {

			// ------------------------------
			// AUDIO WEBM → OGG para WhatsApp
			// ------------------------------
			if input.MessageType == "audio" && input.MIMEType == "audio/webm" {

				fmt.Println("🎧 Recibido audio WebM, convirtiendo a OGG...")

				convertedBase64, newMime, err := convertWebMToOgg(input.Base64Content)
				if err != nil {
					utils.Respond(c, http.StatusInternalServerError, false, "Error convirtiendo audio", nil, err)
					return
				}

				input.Base64Content = convertedBase64
				input.MIMEType = newMime
				input.FileName = "audio.ogg"
			}

			mediaID, err := uploadBase64File(
				*channelIntegration.AppIdentifier,
				*channelIntegration.AccessToken,
				input.Base64Content,
				input.MIMEType,
				input.FileName,
			)
			if err != nil {
				utils.Respond(c, http.StatusInternalServerError, false, "Error al subir audio", nil, err)
				return
			}

			id, err := m.sendWhatsAppAudio(
				*channelIntegration.AppIdentifier,
				*channelIntegration.AccessToken,
				*channelIntegration.SenderID,
				mediaID,
			)
			if err == nil {
				fmt.Println("✅ Audio enviado correctamente. ID:", id)
				input.HasError = false
				input.MessageError = ""
				input.ChannelMessageID = id
			}
			if err != nil {
				utils.Respond(c, http.StatusInternalServerError, false, "Error al enviar audio", nil, err)
				return
			}

		}
	}

	repository := repository.MessageRepository{}

	if err := repository.SendMessageToPlatform(input); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al enviar el mensaje", nil, err)
		return
	}

	utils.InvalidateFirstMessagesByCase(uint(input.CaseID))

	// Consumir

	// payload, _ := json.Marshal(WSMessage{
	// 	Type:   "new_message",
	// 	CaseID: uint(input.CaseID),
	// 	Data:   input, // o arma un DTO si prefieres
	// })
	// channel := "case:" + strconv.Itoa(int(input.CaseID))
	// m.hub.BroadcastJSON(channel, payload)

	utils.Respond(c, http.StatusOK, true, "Mensaje enviado correctamente", input, nil)

	// 🔄 Refresh MV async
	mv := services.NewCaseMVRefreshService()
	mv.RefreshOnEvent("whatsapp_messag_send")
}

func extractBase64(input string) (mime string, raw string) {
	// Esperado: data:<mime>;base64,<contenido>
	if !strings.HasPrefix(input, "data:") {
		return "", input // no trae prefijo, devolver directo
	}

	// Separar MIME y contenido
	parts := strings.SplitN(input, ",", 2)
	if len(parts) != 2 {
		return "", ""
	}

	header := parts[0] // data:application/pdf;base64
	raw = normalizeBase64(parts[1])

	// Extraer MIME
	mime = strings.TrimPrefix(header, "data:")
	mime = strings.TrimSuffix(mime, ";base64")

	return mime, raw
}

func normalizeBase64(b64 string) string {
	missing := len(b64) % 4
	if missing > 0 {
		b64 += strings.Repeat("=", 4-missing)
	}
	return b64
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

func (m *MessageEntry) SendMessengerTextDirect(
	pageAccessToken string,
	recipientPSID string,
	text string,
) (string, error) {

	url := "https://graph.facebook.com/v18.0/me/messages"

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientPSID,
		},
		"message": map[string]string{
			"text": text,
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+pageAccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf(
			"messenger api error (%d): %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	var result struct {
		MessageID string `json:"message_id"`
	}

	_ = json.Unmarshal(respBody, &result)

	return result.MessageID, nil
}

func (m *MessageEntry) SendMessengerImage(
	pageAccessToken string,
	recipientPSID string,
	imageURL string,
) error {

	url := "https://graph.facebook.com/v18.0/me/messages"

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientPSID,
		},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "image",
				"payload": map[string]interface{}{
					"url":         imageURL,
					"is_reusable": false,
				},
			},
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+pageAccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"messenger image error (%d): %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	return nil
}

func (m *MessageEntry) SendMessengerAudio(
	pageAccessToken string,
	recipientPSID string,
	audioURL string,
) error {

	url := "https://graph.facebook.com/v18.0/me/messages"

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientPSID,
		},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "audio",
				"payload": map[string]interface{}{
					"url": audioURL,
				},
			},
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+pageAccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"messenger audio error (%d): %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	return nil
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

// send WhatsApp text direct via API
// func (m *MessageEntry) SendWhatsAppTextDirect(phoneNumberID, accessToken, to, messageText string) error {
// 	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/messages", phoneNumberID)

// 	payload := map[string]interface{}{
// 		"messaging_product": "whatsapp",
// 		"to":                to,
// 		"type":              "text",
// 		"text": map[string]interface{}{
// 			"body": messageText,
// 		},
// 	}

// 	body, _ := json.Marshal(payload)

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return fmt.Errorf("error creando request: %w", err)
// 	}

// 	req.Header.Set("Authorization", "Bearer "+accessToken)
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("error enviando request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	respBody, _ := io.ReadAll(resp.Body)
// 	fmt.Println("📡 Respuesta WhatsApp:", string(respBody))

// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
// 		return fmt.Errorf("error API WhatsApp (%d): %s", resp.StatusCode, string(respBody))
// 	}

// 	return nil
// }

func (m *MessageEntry) SendWhatsAppTextDirect(
	phoneNumberID, accessToken, to, messageText string,
) (string, string, error) {

	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/messages", phoneNumberID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]interface{}{
			"body": messageText,
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", "", fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error enviando request: %w", err)
	}
	defer resp.Body.Close()

	// Respuesta cruda
	respBody, _ := io.ReadAll(resp.Body)
	respString := string(respBody)

	fmt.Println("📡 Respuesta WhatsApp:", respString)

	// Validaciones del status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", respString, fmt.Errorf("error API WhatsApp (%d): %s", resp.StatusCode, respString)
	}

	// Parsear respuesta
	var result struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", respString, fmt.Errorf("error parseando respuesta: %w", err)
	}

	if len(result.Messages) == 0 {
		return "", respString, fmt.Errorf("la API regresó 200 pero no devolvió mensajes")
	}

	wamid := result.Messages[0].ID

	return wamid, respString, nil
}

// Send via n8n WhatsApp webhook
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

func (m *MessageEntry) sendWhatsAppImage(phoneNumberID, accessToken, to, mediaID, caption string) (string, error) {
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
		return "", fmt.Errorf("❌ error creando request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("❌ error enviando request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("📡 Respuesta de envío:", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("❌ error API (%d): %s", resp.StatusCode, string(respBody))
	}

	// Parsear respuesta
	var result struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error parseando respuesta: %w", err)
	}

	if len(result.Messages) == 0 {
		return "", fmt.Errorf("la API regresó 200 pero no devolvió mensajes")
	}

	id := result.Messages[0].ID

	fmt.Println("✅ Mensaje enviado correctamente.")
	return id, nil
}

func textProto(mime, filename string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", mime)
	return h
}

func uploadBase64File(appID, token, rawBase64, mime, fileName string) (string, error) {

	decoded, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		return "", fmt.Errorf("invalid base64: %w", err)
	}

	uploadURL := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/media", appID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// ---- Archivo con MIME real ----
	part, err := writer.CreatePart(textProto(mime, fileName))
	if err != nil {
		return "", err
	}
	part.Write(decoded)

	// ---- Campos obligatorios ----
	writer.WriteField("messaging_product", "whatsapp")
	writer.Close()

	req, _ := http.NewRequest("POST", uploadURL, body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var response struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", err
	}

	if response.ID == "" {
		return "", fmt.Errorf("file upload failed")
	}

	return response.ID, nil
}
func (m *MessageEntry) sendWhatsAppDocument(appID, token, recipient, mediaID, caption string) (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/messages", appID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                recipient,
		"type":              "document",
		"document": map[string]interface{}{
			"id":      mediaID,
			"caption": caption,
		},
	}

	b, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)

	if res.StatusCode >= 300 {
		return "", fmt.Errorf("error sending WhatsApp document: %s", string(respBody))
	}

	// Parsear respuesta
	var result struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error parseando respuesta: %w", err)
	}

	if len(result.Messages) == 0 {
		return "", fmt.Errorf("la API regresó 200 pero no devolvió mensajes")
	}

	return result.Messages[0].ID, nil
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

	// 🔄 Refresh MV async
	mv := services.NewCaseMVRefreshService()
	mv.RefreshOnEvent("whatsapp_message_assign_campaign")
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

	// 🔄 Refresh MV async
	mv := services.NewCaseMVRefreshService()
	mv.RefreshOnEvent("whatsapp_message_assign_department")
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

	// 🔄 Refresh MV async
	mv := services.NewCaseMVRefreshService()
	mv.RefreshOnEvent("whatsapp_message_assign_agent")
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

	// 1. Obtener el caso antes de cerrarlo para saber el agente
	caseData, err := repo.GetCaseByID(uint(req.CaseID))
	if err != nil {
		// Loguear error pero intentar cerrar de todas formas?
		// Mejor fallar si no podemos validar, o simplemente continuar.
		// En este caso continuamos el cierre, el cache expirará solo en 30s si falla esto.
		fmt.Println("⚠️ Error obteniendo caso para invalidar cache:", err)
	}

	if err := repo.CloseCase(req); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al cerrar el caso", nil, err)
		return
	}

	// 2. Invalidar cache si tenía agente asignado
	if caseData != nil && caseData.AgentID > 0 {
		utils.InvalidateActiveCasesByAgent(caseData.AgentID)
		fmt.Println("🧹 Cache invalidado para agente:", caseData.AgentID)
	}

	utils.Respond(c, http.StatusOK, true, "Caso cerrado correctamente", nil, nil)

	// 🔄 Refresh MV async
	mv := services.NewCaseMVRefreshService()
	mv.RefreshOnEvent("whatsapp_message_close_case")
}

//api.POST("/entry/cancel_case/:case_id", controller.CancelCase)

func (m *MessageEntry) CancelCase(c *gin.Context) {
	caseIDStr := c.Param("case_id")
	caseID, err := strconv.Atoi(caseIDStr)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "case_id inválido", nil, err)
		return
	}

	repo := repository.MessageRepository{}
	if err := repo.CancelCase(uint(caseID)); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error al cancelar el caso", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Caso cancelado correctamente", nil, nil)
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

	repo := repository.MessageRepository{}
	cases, err := repo.GetCaseGeneralInformation(uint(companyID), uint(campaignID))
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

func (m *MessageEntry) sendWhatsAppAudio(appID, token, recipient, mediaID string) (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v24.0/%s/messages", appID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                recipient,
		"type":              "audio",
		"audio": map[string]interface{}{
			"id": mediaID,
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error enviando request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("📡 Respuesta WhatsApp AUDIO:", string(respBody))

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("error API (%d): %s", resp.StatusCode, string(respBody))
	}

	// Parsear respuesta
	var result struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error parseando respuesta: %w", err)
	}

	if len(result.Messages) == 0 {
		return "", fmt.Errorf("la API regresó 200 pero no devolvió mensajes")
	}

	return result.Messages[0].ID, nil
}

func cleanBase64(input string) string {
	// Remove prefix like: data:audio/webm;base64,
	if strings.Contains(input, ",") {
		parts := strings.SplitN(input, ",", 2)
		input = parts[1]
	}

	// Remove invalid characters
	re := regexp.MustCompile(`[^A-Za-z0-9+/=]`)
	input = re.ReplaceAllString(input, "")

	// Fix missing padding
	if m := len(input) % 4; m != 0 {
		input += strings.Repeat("=", 4-m)
	}

	return input
}

// func convertWebMToOgg(b64 string) (string, string, error) {

// 	// NORMALIZAR BASE64
// 	clean := strings.ReplaceAll(b64, "\n", "")
// 	clean = strings.ReplaceAll(clean, "\r", "")
// 	clean = strings.ReplaceAll(clean, " ", "")
// 	missing := len(clean) % 4
// 	if missing > 0 {
// 		clean += strings.Repeat("=", 4-missing)
// 	}

// 	decoded, err := base64.StdEncoding.DecodeString(clean)
// 	if err != nil {
// 		return "", "", fmt.Errorf("decode error: %w", err)
// 	}

// 	// Guardar archivo temporal
// 	if err := os.WriteFile("tmp_in.webm", decoded, 0644); err != nil {
// 		return "", "", err
// 	}

// 	// Forzar formato webm →
// 	cmd := exec.Command("ffmpeg",
// 		"-y",
// 		"-f", "webm",
// 		"-i", "tmp_in.webm",
// 		"-c:a", "libopus",
// 		"tmp_out.ogg",
// 	)

// 	var stderr bytes.Buffer
// 	cmd.Stderr = &stderr

// 	if err := cmd.Run(); err != nil {
// 		return "", "", fmt.Errorf("ffmpeg error: %v | stderr: %s", err, stderr.String())
// 	}

// 	// Leer ogg resultante
// 	oggBytes, err := os.ReadFile("tmp_out.ogg")
// 	if err != nil {
// 		return "", "", err
// 	}

// 	oggBase64 := base64.StdEncoding.EncodeToString(oggBytes)

// 	return oggBase64, "audio/ogg", nil
// }

func convertWebMToOgg(b64 string) (string, string, error) {

	// 1️⃣ Normalizar base64
	clean := strings.ReplaceAll(b64, "\n", "")
	clean = strings.ReplaceAll(clean, "\r", "")
	clean = strings.ReplaceAll(clean, " ", "")
	if m := len(clean) % 4; m != 0 {
		clean += strings.Repeat("=", 4-m)
	}

	decoded, err := base64.StdEncoding.DecodeString(clean)
	if err != nil {
		return "", "", fmt.Errorf("decode error: %w", err)
	}

	// 2️⃣ Crear carpeta temporal única
	tempDir, err := os.MkdirTemp("", "audio_convert_*")
	if err != nil {
		return "", "", fmt.Errorf("temp dir error: %w", err)
	}

	// borrar todo al terminar
	defer os.RemoveAll(tempDir)

	// 3️⃣ Rutas seguras
	inPath := filepath.Join(tempDir, "in.webm")
	outPath := filepath.Join(tempDir, "out.ogg")

	// 4️⃣ Guardar el archivo webm
	if err := os.WriteFile(inPath, decoded, 0644); err != nil {
		return "", "", fmt.Errorf("write error: %w", err)
	}

	// 5️⃣ Ejecutar ffmpeg con formato forzado
	cmd := exec.Command("ffmpeg",
		"-y",
		"-f", "webm",
		"-i", inPath,
		"-c:a", "libopus",
		outPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("ffmpeg error: %v | stderr: %s", err, stderr.String())
	}

	// 6️⃣ Leer ogg final
	oggBytes, err := os.ReadFile(outPath)
	if err != nil {
		return "", "", fmt.Errorf("read ogg error: %w", err)
	}

	oggBase64 := base64.StdEncoding.EncodeToString(oggBytes)

	return oggBase64, "audio/ogg", nil
}

func (m *MessageEntry) DownloadMessageFile(c *gin.Context) {
	messageID := c.Param("message_id")
	if messageID == "" {
		utils.Respond(c, http.StatusBadRequest, false, "message_id requerido", nil, nil)
		return
	}

	id, err := strconv.Atoi(messageID)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "message_id inválido", nil, err)
		return
	}

	repo := repository.MessageRepository{}
	message, err := repo.GetMessageByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, false, "Mensaje no encontrado", nil, err)
		return
	}

	if message.Base64Content == "" {
		utils.Respond(c, http.StatusBadRequest, false, "El mensaje no contiene archivo", nil, nil)
		return
	}

	// Limpiar prefijo data:xxx;base64, si existe
	b64 := message.Base64Content
	if idx := strings.Index(b64, ","); idx != -1 {
		b64 = b64[idx+1:]
	}

	// Decodificar base64
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error decodificando archivo", nil, err)
		return
	}

	// Nombre del archivo
	fileName := message.TextContent
	if fileName == "" {
		// Fallback simple si no hay nombre
		fileName = fmt.Sprintf("file_%d", message.ID)

		// Intentar adivinar extensión por mime
		if strings.Contains(message.MIMEType, "pdf") {
			fileName += ".pdf"
		} else if strings.Contains(message.MIMEType, "image") {
			fileName += ".jpg" // genérico
		} else if strings.Contains(message.MIMEType, "audio") {
			fileName += ".ogg" // por defecto de whatsapp
		} else {
			fileName += ".bin"
		}
	}

	// Limpiar nombre de archivo para seguridad básica
	fileName = filepath.Base(fileName)

	// Directorio de destino
	baseDir := "public/downloads"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error creando directorio de descargas", nil, err)
		return
	}

	fullPath := filepath.Join(baseDir, fileName)

	// Escribir archivo
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error escribiendo archivo", nil, err)
		return
	}

	// Construir URL pública relativa o absoluta según necesidad
	// El usuario pidió "servir como ruta". Devolvemos la ruta relativa pública.
	// Asumimos que ./public es root o similar
	publicURL := fmt.Sprintf("/public/downloads/%s", fileName)

	utils.Respond(c, http.StatusOK, true, "Archivo generado correctamente", gin.H{
		"url":  publicURL,
		"path": fullPath,
	}, nil)
}
