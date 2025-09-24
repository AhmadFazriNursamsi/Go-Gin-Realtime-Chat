package routes

import (
	"myapi/controllers"
	"myapi/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {

	usergroup := r.Group("/users")
	usergroup.Use(middlewares.AuthMiddleware())
	usergroup.Use(middlewares.RequirePermission(db, "permission.manage")) // Hanya admin yang bisa manage user
	{
		usergroup.GET("/", controllers.GetUsers)
		usergroup.POST("/", controllers.CreateUser)
		usergroup.PUT("/:id", controllers.UpdateUser)
		usergroup.DELETE("/:id", controllers.DeleteUser)
		usergroup.POST("/upload/:id", controllers.UploadUserPhoto)
		usergroup.GET("/with-permissions", controllers.GetUsersWithPermissions(db))

	}
}

// func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
// 	userGroup := r.Group("/users")
// 	userGroup.Use(middlewares.AuthMiddleware())
// 	{
// 		// route default user
// 		userGroup.GET("/", controllers.GetUsers(db))

// 		// ✅ route baru untuk ambil user + permissions
// 		userGroup.GET("/with-permissions", controllers.GetUsersWithPermissions(db))
// 	}
// }
