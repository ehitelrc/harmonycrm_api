package controllers

import (
	"encoding/json"
	"harmony_api/ws"

	"github.com/gin-gonic/gin"
)

type WebSocketTestController struct {
	// Add any required fields here
	hub *ws.Hub
}

func NewWebSocketTestController(h *ws.Hub) *WebSocketTestController {
	return &WebSocketTestController{hub: h}
}

func (m WebSocketTestController) SendNotificationByAgent(c *gin.Context) {
	agentID := c.Param("id")

	payload, _ := json.Marshal(WSMessage{
		Type:   "new_message",
		CaseID: 123,          // ejemplo de case_id
		Data:   "newMessage", // o arma un DTO si prefieres
	})
	channel := "case:" + "123" // Aquí deberías obtener el case_id real asociado al agent_id
	m.hub.BroadcastJSON(channel, payload)

	c.JSON(200, gin.H{
		"success":  true,
		"message":  "Notificación enviada al agente",
		"agent_id": agentID,
	})

}

func (ctrl *WebSocketTestController) SendNotificationByCaseID(c *gin.Context) {
	caseID := c.Param("id")

	payload, _ := json.Marshal(WSMessage{
		Type:   "new_message",
		CaseID: 123,          // ejemplo de case_id
		Data:   "newMessage", // o arma un DTO si prefieres
	})
	channel := "case:" + caseID
	ctrl.hub.BroadcastJSON(channel, payload)

	c.JSON(200, gin.H{
		"success": true,
		"message": "Notificación enviada al case",
		"case_id": caseID,
	})
}
