package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterTemplateRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewTemplateController()

	templates := r.Group("/templates")
	templates.GET("/", ctrl.GetAllTemplates)
	templates.GET("/:id", ctrl.GetTemplateByID)
	templates.POST("/", ctrl.CreateTemplate)
	templates.PUT("/:id", ctrl.UpdateTemplate)
	templates.DELETE("/:id", ctrl.DeleteTemplate)
	templates.POST("/:id/register-meta", ctrl.RegisterMetaTemplate)
	templates.POST("/:id/sync-meta", ctrl.SyncMetaTemplate)

	// Integrations linked to a specific template (is_linked view)
	templates.GET("/:id/integrations", ctrl.GetIntegrationsForTemplate)

	//  Relation between templates and integrations
	templates.GET("/integration/:id", ctrl.GetTemplatesByIntegration)
	templates.POST("/integration/:id", ctrl.CreateTemplateIntegration)
	templates.DELETE("/integration/:id", ctrl.DeleteTemplateIntegration)

	// Fetch template preview from Meta
	templates.GET("/preview/:template_name/:integration_id", ctrl.PreviewMetaTemplate)
}
