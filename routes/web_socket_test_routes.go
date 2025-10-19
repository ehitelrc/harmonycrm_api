package routes

import (
	"harmony_api/controllers"
	"harmony_api/ws"

	"github.com/gin-gonic/gin"
)

func RegisterWebSocketTestRoutes(r *gin.RouterGroup, hub *ws.Hub) {
	ctrl := controllers.NewWebSocketTestController(hub)

	r.POST("/web-socket-test/agent/:id", ctrl.SendNotificationByAgent)
	r.GET("/web-socket-test/case/:id", ctrl.SendNotificationByCaseID)
}
