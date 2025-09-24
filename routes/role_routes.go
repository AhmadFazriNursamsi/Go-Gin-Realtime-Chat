package routes

import (
	"myapi/controllers"
	"myapi/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoleRoutes(r *gin.Engine, db *gorm.DB) {
	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization
	rolegroup := r.Group("/roles")
	rolegroup.Use(middlewares.AuthMiddleware())
	rolegroup.Use(middlewares.RequireRole(2)) // Hanya role dengan ID 1 (Admin) yang boleh akses
	{
		rolegroup.GET("/", controllers.GetRoles)
		rolegroup.POST("/", controllers.CreateRole)
		rolegroup.PUT("/:id", controllers.UpdateRole)
		rolegroup.DELETE("/:id", controllers.DeleteRole)
	}
}
