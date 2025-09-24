package middlewares

import (
	// "log"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequirePermission(db *gorm.DB, requiredPerm string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleID := ctx.GetInt("role_id")
		if roleID == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Missing role_id in token"})
			return
		}

		// Ambil permission dari DB
		var count int64
		err := db.Table("role_permissions").
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Where("role_permissions.role_id = ? AND permissions.name = ?", roleID, requiredPerm).
			Count(&count).Error

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
			return
		}

		log.Printf("[DEBUG] Checking permission: role_id=%d, required=%s", roleID, requiredPerm)
		if count == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		ctx.Next()
	}
}
