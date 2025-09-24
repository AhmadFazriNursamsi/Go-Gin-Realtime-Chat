package routes

import (
	"myapi/controllers"
	"myapi/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterPermissionRoutes(r *gin.Engine, db *gorm.DB) {
	permGroup := r.Group("/permissions")
	permGroup.Use(middlewares.AuthMiddleware())
	permGroup.Use(middlewares.RequirePermission(db, "permission.manage")) // Hanya admin yang bisa manage permission
	{

		permGroup.GET("/", controllers.GetPermissions(db))
		permGroup.POST("/", controllers.CreatePermission(db))
		permGroup.DELETE("/:id", controllers.DeletePermission(db))
	}
}

func RegisterRolePermissionRoutes(r *gin.Engine, db *gorm.DB) {
	rpGroup := r.Group("/role-permissions")
	rpGroup.Use(middlewares.AuthMiddleware())
	rpGroup.Use(middlewares.RequirePermission(db, "permission.manage"))
	{
		rpGroup.POST("/", controllers.AssignPermissionToRole(db))
		rpGroup.DELETE("/:roleID/:permissionID", controllers.RemovePermissionFromRole(db))
	}
}
