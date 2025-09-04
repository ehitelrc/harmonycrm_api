package routes

import (
	"github.com/gin-gonic/gin"

	"harmony_api/controllers"
)

func RegisterFunnelRoutes(r *gin.RouterGroup) {
	ctl := controllers.NewFunnelController()

	g := r.Group("/funnels")
	{
		g.GET("", ctl.GetAll)
		g.GET("/:id", ctl.GetByID)
		g.POST("/", ctl.Create)
		g.PUT("/", ctl.Update)
		g.DELETE("/:id", ctl.Delete)

		// Stages

		g.GET("/:id/stages", ctl.GetStages)
		g.GET("/:id/stages/:stage_id", ctl.GetStageByID)
		g.POST("/stages", ctl.CreateStage)
		g.PUT("/stages", ctl.UpdateStage)
		g.DELETE("/stages/:stage_id", ctl.DeleteStage)
	}
}
