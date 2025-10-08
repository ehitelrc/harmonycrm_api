package routes

import (
	"harmony_api/controllers"
	"harmony_api/ws"

	"github.com/gin-gonic/gin"
)

func InitializeMessage(r gin.RouterGroup, hub *ws.Hub) {

	controller := controllers.NewMessageEntry(hub)

	api := r.Group("/messages")

	api.POST("/entry", controller.ReceiveMessageWebhook)

	api.POST("/entry/ws/media/image", controller.ReceiveImageMessageWebhookMedia)

	// Endpoint para recibir mensajes de audio
	api.POST("/entry/ws/media/audio", controller.ReceiveAudioMessageWebhookMedia)

	// Active cases by agent_id
	api.GET("/entry/active_cases/:agent_id", controller.GetActiveCasesByAgentID)

	// Get messages by case_id
	api.GET("/entry/messages/:case_id", controller.GetMessagesByCaseID)

	// Send message to platform
	api.POST("/entry/send", controller.SendMessageToPlatform)

	// Assign case to client
	api.PUT("/entry/assign_case", controller.AssignCaseToClient)

	//Add case notes
	api.POST("/entry/case_notes", controller.AddCaseNote)

	// Get notes by case_id
	api.GET("/entry/case_notes/:case_id", controller.GetCaseNotesByCaseID)

	// Assign to campaign
	api.POST("/entry/assign_campaign", controller.AssignCaseToCampaign)

	api.POST("/entry/assign_department", controller.AssignCaseToDepartment)

	// Assign to agent
	api.POST("/entry/assign_agent", controller.AssignCaseToAgent)

	// Current case funnel
	api.GET("/entry/case_funnel/current/:case_id", controller.GetCurrentCaseFunnel)

	// Set case funnel stage
	api.POST("/entry/case_funnel/set_stage", controller.SetCaseFunnelStage)

	// Close case
	api.POST("/entry/close_case", controller.CloseCase)

	// Get case general information with company_id, campaign_id and stage_id
	api.GET("/entry/case_general_info/:company_id/:campaign_id/:stage_id", controller.GetCaseGeneralInformation)

	// Get unassigned cases by company_id
	api.GET("/entry/unassigned_cases/:company_id", controller.GetCasesWithoutAgentByCompanyID)
}
