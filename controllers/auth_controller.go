package controllers

import (
	"log"
	"myapi/models"
	"myapi/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// var jwtKey []byte
var jwtSecret []byte // ✅ ini jadi source of truth

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found")
	}

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Println("⚠️  JWT_SECRET tidak ditemukan di env, gunakan default secret.")
		jwtSecret = []byte("0WvRtY6h9V7qCrMm6KDxjD3c6nQFlQ0gTK9r4ggh7LM=")
	}
	log.Printf("🔑 JWT secret loaded, length=%d\n", len(jwtSecret))

}

type LoginInput struct {
	Username string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type token struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}
type ErrorResponse struct {
	Message string `json:"message" example:"Email atau username tidak ditemukan"`
}

// LoginHandler godoc
// @Summary Login user
// @Description Login untuk mendapatkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login"
// @Success 200 {object} token
// @Failure 400 {object} ErrorResponse
// @Router /login [post]
func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input LoginInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.Preload("Role").
			Preload("RoleChild").
			Preload("Rooms").
			Where("email = ?", input.Username).
			First(&user).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau username tidak ditemukan"})
			return
		}

		if !utils.CheckPasswordHash(input.Password, user.Password) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "password salah"})
			return
		}

		claims := models.CustomClaims{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
			Roleid:        user.RoleID,
			RoleName:      user.Role.Name, // nanti diisi kalau perlu
			RoleChildID:   user.RoleChildID,
			RoleChildName: user.RoleChild.Name, // pastikan RoleChild sudah di-preload
			// RoomsId:       roomID,

			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				Issuer:    "Atep Token",
				Subject:   strconv.Itoa(int(user.ID)),
			},
		}
		var roomIDs []uint
		for _, r := range user.Rooms {
			roomIDs = append(roomIDs, r.ID)
		}
		claims.RoomsId = roomIDs

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// tokenString, err := token.SignedString(jwtKey)
		tokenString, err := token.SignedString(jwtSecret)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		log.Printf("[DEBUG] JWT_SECRET length: %d\n", len(jwtSecret))
		log.Println("Loaded JWT secret:", string(jwtSecret))
		log.Println("Length:", len(jwtSecret))
		ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}
func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input struct {
			Name        string `json:"name" binding:"required"`
			Email       string `json:"email" binding:"required,email"`
			Password    string `json:"password" binding:"required"`
			Roleid      *uint  `json:"role_id"`
			RoleChildID *uint  `json:"rolechild_id"`
		}
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var existing models.User
		if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
			return
		}
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
			return
		}
		user := models.User{
			Name:        input.Name,
			Email:       input.Email,
			Password:    hashedPassword,
			RoleID:      input.Roleid,
			RoleChildID: input.RoleChildID,
		}
		if err := db.Create(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}

}
func ProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetInt("user_id")
		role := ctx.GetString("role_name") // ✅ bisa akses role juga kalau perlu

		if id == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing user_id"})
			return
		}

		var user models.User
		if err := db.Preload("RoleChild").First(&user, id).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Kamu bisa kirim data role juga kalau mau
		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
			"role": role,
		})
	}
}
func LogoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Karena JWT bersifat stateless, logout biasanya dilakukan di sisi klien
		ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful on client side"})
	}
}
func ForgotPasswordHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input struct {
			Email string `json:"email" binding:"required,email"`
		}
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}
		// Di sini kamu bisa generate token reset dan kirim email
		ctx.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to email (not implemented)"})
	}
}

func ResetPasswordHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input struct {
			Email       string `json:"email" binding:"required,email"`
			NewPassword string `json:"new_password" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}
		hashedPassword, err := utils.HashPassword(input.NewPassword)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
			return
		}
		user.Password = hashedPassword
		db.Save(&user)
		ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
	}
}
