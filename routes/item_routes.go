package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterItemRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewItemController()

	// CRUD básico
	r.GET("/items", ctrl.GetAll)
	r.GET("/items/:id", ctrl.GetByID)

	// Filtro por compañía
	r.GET("/items/company/:company_id", ctrl.GetByCompany)

	// Crear / Actualizar (PUT con objeto completo que incluya id)
	r.POST("/items", ctrl.Create)
	r.PUT("/items", ctrl.Update)

	// Eliminar
	r.DELETE("/items/:id", ctrl.Delete)
}
