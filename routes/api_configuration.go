package routes

import (
	"fmt"
	"harmony_api/config"
	"harmony_api/controllers"
	"harmony_api/providers"
	"harmony_api/repository"
	"harmony_api/services"
	"harmony_api/ws"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine, hub *ws.Hub) {

	// Crear OpenAIProvider
	openAIProvider, err := providers.NewOpenAIProvider("gpt-4o-mini")
	if err != nil {
		panic(err)
	}

	googleOCR, errOcr := providers.NewGoogleOCR("IA/harmonyvpocr-07dfbad813ad.json")
	if errOcr != nil {
		fmt.Printf("❌ [OCR Init] Error initializing Google OCR: %v\n", errOcr)
	} else {
		fmt.Println("✅ [OCR Init] Google OCR initialized successfully!")
	}
	ocrService := services.NewOCRService(googleOCR)
	ocrController := controllers.NewOCRController(ocrService)

	// Receipt State Management
	receiptRepo := repository.NewReceiptRepository()
	receiptStateService := services.NewReceiptStateService(receiptRepo)
	receiptStateController := controllers.NewReceiptStateController(receiptStateService)

	// Crear ReceiptAnalysisService
	receiptAnalysisService := services.NewReceiptAnalysisService(ocrService, openAIProvider, repository.NewReceiptRepository())

	// Crear controller
	receiptAnalysisController := controllers.NewReceiptAnalysisController(receiptAnalysisService)

	// Obtener el path absoluto desde la raíz del proyecto (subiendo desde cmd/)
	rootDir, _ := filepath.Abs(filepath.Join(".", ".."))
	assetsPath := filepath.Join(rootDir, "assets")
	r.Static("/assets", assetsPath)

	api := r.Group("/api")

	// Inicializar rutas de mensajes
	InitializeMessage(*api, hub, receiptAnalysisService)

	// Inicializar rutas de casos cerrados por Sender ID
	messageRepo := repository.MessageRepository{}
	closedCasesController := controllers.NewClosedCasesController(&messageRepo)
	RegisterClosedCasesRoutes(api, closedCasesController)

	// Inicializar rutas de compañías
	RegisterCompanyRoutes(api)

	// Inicializar rutas de departamentos
	RegisterDepartmentRoutes(api)

	// Inicializar rutas de canales
	RegisterChannelRoutes(api)

	// Inicializar rutas de campañas
	RegisterCampaignRoutes(api)

	// Inicializar rutas de clientes
	RegisterClientRoutes(api)

	// Inicializar rutas de agentes
	RegisterAgentRoutes(api)

	// Inicializar rutas de usuarios
	RegisterUserRoutes(api)

	// Inicializar rutas de configuración
	RegisterAgentDepartmentAssignmentRoutes(api)

	// Inicializar rutas de roles
	RegisterRoleRoutes(api)

	// Inicializar rutas de permisos
	RegisterPermissionRoutes(api)

	// Inicializar rutas de asignación de permisos a roles
	RegisterRolePermissionRoutes(api)

	// Inicializar rutas de asignación de roles a usuarios en compañías
	RegisterUserCompanyRoleRoutes(api)

	// Inicializar rutas de cuentas sociales de clientes
	RegisterClientSocialAccountRoutes(api)

	// Login
	RegisterLoginRoutes(api)

	// Inicializar rutas de items
	RegisterItemRoutes(api)

	// Inicializar rutas de embudos
	RegisterFunnelRoutes(api)

	// Dashboard
	RegisterDashboardRoutes(api)

	// Campaign Pushing
	RegisterCampaignPushingRoutes(api, hub)

	// Locations
	LocationsRoutes(api)

	RegisterCaseItemRoutes(api)

	RegisterCaseDashboardRoutes(api)

	RegisterCasesBulkRoutes(api)

	RegisterWebSocketTestRoutes(api, hub)

	RegisterCustomListRoutes(api)

	RegristerCustomFieldRoutes(api)

	RegisterWhatsappWebHookRoutes(api, hub, receiptAnalysisService)

	OCRRoutes(api, ocrController)

	// Registrar rutas
	ReceiptAnalysisRoutes(api, receiptAnalysisController)

	ReceiptStateRoutes(api, receiptStateController)

	RegisterMessengerWebHookRoutes(api, hub, receiptAnalysisService)

	RegisterTemplateRoutes(api)
	RegisterTemplateReportRoutes(api)
	RegisterOcrReportRoutes(api)
	RegisterMessageStatusReportRoutes(api)

	RegisterPaymentValidationRoutes(api)

	RegisterUnreconciledPaymentsRoutes(api)

	RegisterTagRoutes(api)

	r.Static("/public", "./public")

	// Endpoint de verificación de estado
	api.GET("/health", func(c *gin.Context) {
		version := "3.50.0"
		if versionBytes, err := os.ReadFile("version.txt"); err == nil {
			version = strings.TrimSpace(string(versionBytes))
		}

		dbStatus := "Disconnected"
		// Asumiendo que config.DB es la instancia *gorm.DB global
		if sqlDB, err := config.DB.DB(); err == nil {
			if err := sqlDB.Ping(); err == nil {
				dbStatus = "Connected"
			} else {
				dbStatus = "Error: " + err.Error()
			}
		}

		c.JSON(200, gin.H{
			"status":   "API is online",
			"version":  version,
			"database": dbStatus,
		})
	})

}
