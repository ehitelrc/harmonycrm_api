package routes

import (
	"harmony_api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRolePermissionRoutes(r *gin.RouterGroup) {
	ctrl := controllers.NewRolePermissionController()

	// Consultas
	r.GET("/role-permissions/role/:role_id", ctrl.GetByRole)
	r.GET("/role-permissions/permission/:permission_id", ctrl.GetByPermission)

	// Asignar / Desasignar uno
	r.POST("/role-permissions", ctrl.Assign)
	r.DELETE("/role-permissions/role/:role_id/permission/:permission_id", ctrl.Unassign)

	// Reemplazar todos los permisos de un rol (transacción)
	r.PUT("/role-permissions/role/:role_id", ctrl.ReplaceForRole)
}
