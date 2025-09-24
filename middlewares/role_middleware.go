package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the user has the required role to access the endpoint
// ...existing code...

// RequireRole checks if the user has the required role to access the endpoint
func RequireRole(requiredRoles ...int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		RoleID := ctx.GetInt("role_id")
		log.Println("🔎 RoleID dari token:", RoleID) // DEBUG
		if RoleID == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found in token"})
			return
		}
		allowed := false
		for _, role := range requiredRoles {
			if role == RoleID {
				allowed = true
				break
			}
		}
		if !allowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: insufficient role"})
			return
		}
		ctx.Next()
	}
}
