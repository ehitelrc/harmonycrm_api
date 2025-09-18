package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAgentRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewAgentController()

	r.GET("/agents", ctrl.GetAll)
	r.GET("/agents/:user_id", ctrl.GetByUserID)
	r.POST("/agents", ctrl.Create)
	r.DELETE("/agents/:user_id", ctrl.Delete)

	// Agentes con info de usuario
	r.GET("/agents/agents-with-user-info", ctrl.GetAllWithUserInfo)
	r.GET("/agents/agents-with-user-info/:user_id", ctrl.GetByUserIDWithUserInfo)

	// Non agents users
	r.GET("/agents/non-agents", ctrl.GetAllNonAgents)
}
