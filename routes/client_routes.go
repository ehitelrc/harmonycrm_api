package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewClientController()

	r.GET("/clients", ctrl.GetAll)
	r.GET("/clients/:id", ctrl.GetByID)
	r.POST("/clients", ctrl.Create)
	r.PUT("/clients", ctrl.Update) // objeto completo con id
	r.DELETE("/clients/:id", ctrl.Delete)

	// Create a new lead
	r.POST("/clients/leads", ctrl.CreateLead)

	r.GET("/clients/custom_fields/:entity_id", ctrl.GetCustomFields)

	// ✅ Duplicados por teléfono
	r.GET("/clients/duplicates/phone", ctrl.GetDuplicatePhonesDTO)

	// routes/client_routes.go
	r.GET("/clients/external/:external_id", ctrl.GetByExternalID)
}
