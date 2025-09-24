package middlewares

import (
	"errors"
	"myapi/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey []byte

// Profile godoc
// @Summary      Profile user
// @Description  get user info
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.User
// @Failure      401  {object}  models.User
// @Failure      403  {object}  models.User
// @Router /users/with-permissions [get]
// @Security Bearer
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &models.CustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			if errors.Is(err, jwt.ErrTokenExpired) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		uid, err := strconv.Atoi(claims.Subject)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		ctx.Set("user_id", uid)

		// ✅ Pastikan role_id selalu ada (0 jika nil)
		if claims.Roleid != nil {
			ctx.Set("role_id", int(*claims.Roleid))
		} else {
			ctx.Set("role_id", 0)
		}

		ctx.Set("role_name", claims.RoleName)

		ctx.Next()
	}
}
