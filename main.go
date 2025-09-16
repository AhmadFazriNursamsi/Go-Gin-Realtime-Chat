package main

import (
	"myapi/controllers"
	"myapi/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()

	usergroup := r.Group("/users")
	{
		usergroup.GET("/", controllers.GetUsers)
		// usergroup.GET("/:id", controllers.GetUser)
		usergroup.POST("/", controllers.CreateUser)
		usergroup.PUT("/:id", controllers.UpdateUser)
		usergroup.DELETE("/:id", controllers.DeleteUser)
		usergroup.POST("/upload/:id", controllers.UploadUserPhoto)
	}
	rolegroup := r.Group("/roles")
	{
		rolegroup.GET("/", controllers.GetRoles)
		rolegroup.POST("/", controllers.CreateRole)
		rolegroup.PUT("/:id", controllers.UpdateRole)
		// rolegroup.POST("/assign/:userID/:roleID", controllers.AssignRole)
		rolegroup.DELETE("/remove/:userID/:roleID", controllers.DeleteRole)
	}
	rolechildgroup := r.Group("/rolechild")
	{
		rolechildgroup.GET("/", controllers.Getrolechild)
		rolechildgroup.POST("/", controllers.Createrolechild)
		rolechildgroup.PUT("/:id", controllers.Updaterolechild)
		rolechildgroup.DELETE("/:id", controllers.Deleterolechild)
	}

	r.Run(":8080")
}
