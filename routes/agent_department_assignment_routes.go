package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAgentDepartmentAssignmentRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewAgentDepartmentAssignmentController()

	r.GET("/agent-department-assignments/:id", ctrl.GetByID)
	r.GET("/agent-department-assignments/agent/:agent_id", ctrl.GetByAgent)
	r.GET("/agent-department-assignments/department/:department_id", ctrl.GetByDepartment)
	r.POST("/agent-department-assignments", ctrl.Create)
	r.PUT("/agent-department-assignments", ctrl.Update) // objeto completo con id
	r.DELETE("/agent-department-assignments/:id", ctrl.Delete)
}
