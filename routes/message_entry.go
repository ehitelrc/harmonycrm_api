package routes

import (
	"harmony_api/controllers"
	"harmony_api/services"
	"harmony_api/ws"

	"github.com/gin-gonic/gin"
)

func InitializeMessage(r gin.RouterGroup, hub *ws.Hub, receiptAnalysisService *services.ReceiptAnalysisService) {

	controller := controllers.NewMessageEntry(hub, receiptAnalysisService)

	api := r.Group("/messages")

	api.POST("/entry", controller.ReceiveMessageWebhook)

	api.POST("/entry/ws/media/image", controller.ReceiveImageMessageWebhookMedia)

	// Endpoint para recibir mensajes de audio
	api.POST("/entry/ws/media/audio", controller.ReceiveAudioMessageWebhookMedia)

	// Active cases by agent_id
	api.GET("/entry/active_cases/:agent_id", controller.GetActiveCasesByAgentID)

	// Get messages by case_id
	api.GET("/entry/messages/:case_id", controller.GetMessagesByCaseID)

	// Download message file
	api.GET("/entry/download/:message_id", controller.DownloadMessageFile)

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

	// Cancel case
	api.POST("/entry/cancel_case/:case_id", controller.CancelCase)

	// Get case general information with company_id, campaign_id
	api.GET("/entry/case_general_info/:company_id/:campaign_id", controller.GetCaseGeneralInformation)

	// vw_case_general_information Get by company_id, campaign_id and agent_id
	api.GET("/entry/leads/company/:company_id/campaign/:campaign_id/agent/:agent_id/channel_integration/:channel_integration_id", controller.GetCaseGeneralInformationByCompanyCampaignAgent)

	// Get unassigned cases by company_id
	api.GET("/entry/unassigned_cases/:company_id", controller.GetCasesWithoutAgentByCompanyID)

	// By company and department unassigned cases
	api.GET("/entry/unassigned_cases/company/:company_id/department/:department_id", controller.GetCasesWithoutAgentByCompanyAndDepartmentID)

	// Open cases by company_id and department_id
	api.GET("/entry/open_cases/company/:company_id/department/:department_id", controller.GetOpenCasesByCompanyAndDepartmentID)

	api.GET("/v2/entry/open_cases/company/:company_id/department/:department_id", controller.GetOpenCasesByCompanyAndDepartmentIDV2)

	api.GET(
		"/v2/entry/open_cases/stats/company/:company_id/department/:department_id",
		controller.GetOpenCasesStatsByCompanyAndDepartmentV2,
	)

	api.GET(
		"/v2/entry/open_cases_mv/company/:company_id/department/:department_id",
		controller.GetOpenCasesMV,
	)

	api.GET(
		"/v2/entry/stats/company/:company_id/department/:department_id",
		controller.GetCaseStats,
	)

	//

	// Mark messages as read by case_id
	api.PUT("/entry/mark_messages_read/case/:case_id", controller.MarkMessagesAsReadByCaseID)

	// New case from template
	api.POST("/entry/new-case/template", controller.CreateNewCaseFromTemplate)

}
