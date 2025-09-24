package routes

import (
	"myapi/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoleChildRoutes(r *gin.Engine, db *gorm.DB) {
	rolechildgroup := r.Group("/rolechild")
	{
		rolechildgroup.GET("/", controllers.Getrolechild)
		rolechildgroup.POST("/", controllers.Createrolechild)
		rolechildgroup.PUT("/:id", controllers.Updaterolechild)
		rolechildgroup.DELETE("/:id", controllers.Deleterolechild)
	}
	// rolechildgroup.Use(middlewares.AuthMiddleware())
	// {
	// 	rolechildgroup.POST("/assign/:roleID/:childRoleID", controllers.AssignChildRole)
	// }
}
