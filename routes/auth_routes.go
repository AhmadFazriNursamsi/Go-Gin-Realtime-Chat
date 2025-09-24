package routes

import (
	"myapi/controllers"
	"myapi/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB, jwtKey []byte) {
	middlewares.JwtKey = jwtKey // share jwt key ke middleware

	r.POST("/login", controllers.LoginHandler(db))
	r.POST("/register", controllers.RegisterHandler(db))
	r.POST("/logout", controllers.LogoutHandler())
	r.POST("/forgot-password", controllers.ForgotPasswordHandler(db))
	r.POST("/reset-password", controllers.ResetPasswordHandler(db))

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	protected.GET("/profile", controllers.ProfileHandler(db))
}
