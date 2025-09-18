package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserCompanyRoleRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewUserCompanyRoleController()

	// CRUD
	r.GET("/user-company-roles", ctrl.GetAll)
	r.GET("/user-company-roles/:id", ctrl.GetByID)
	r.POST("/user-company-roles", ctrl.Create)
	r.PUT("/user-company-roles", ctrl.Update) // objeto completo con id
	r.DELETE("/user-company-roles/:id", ctrl.Delete)

	// Filtros existentes
	r.GET("/user-company-roles/user/:user_id", ctrl.GetByUser)
	r.GET("/user-company-roles/company/:company_id", ctrl.GetByCompany)
	r.GET("/user-company-roles/role/:role_id", ctrl.GetByRole)

	// 🔥 Nuevos endpoints solicitados
	// 1) Usuarios y permisos por compañía y rol
	r.GET("/user-company-roles/company/:company_id/role/:role_id/users-permissions", ctrl.GetUsersAndPermissionsByCompanyRole)

	// 2) Permisos por compañía y usuario (efectivos por sus roles en esa company)
	r.GET("/user-company-roles/company/:company_id/user/:user_id/permissions", ctrl.GetPermissionsByCompanyUser)

	r.GET("/user-company-roles/user/:user_id/company/:company_id", ctrl.GetByUserAndCompanyMixed)

	// Batch update
	r.PUT("/user-company-roles/batch", ctrl.BatchUpdate)
}
