package routes

import (
	"harmony_api/ws"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine, hub *ws.Hub) {

	// Obtener el path absoluto desde la raíz del proyecto (subiendo desde cmd/)
	rootDir, _ := filepath.Abs(filepath.Join(".", ".."))
	assetsPath := filepath.Join(rootDir, "assets")
	r.Static("/assets", assetsPath)

	api := r.Group("/api")

	// Inicializar rutas de mensajes
	InitializeMessage(*api, hub)

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
	RegisterCampaignPushingRoutes(api)

	// Locations
	LocationsRoutes(api)

	RegisterCaseItemRoutes(api)

	// Endpoint de verificación de estado
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "API is online",
		})
	})

}
