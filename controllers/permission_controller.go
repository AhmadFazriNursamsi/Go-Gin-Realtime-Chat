package controllers

import (
	"myapi/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPermissions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var permissions []models.Permission
		if err := db.Find(&permissions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get permissions"})
			return
		}
		c.JSON(http.StatusOK, permissions)
	}
}

func CreatePermission(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		perm := models.Permission{Name: input.Name}
		if err := db.Create(&perm).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permission"})
			return
		}
		ctx.JSON(http.StatusOK, perm)
	}
}

func DeletePermission(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := db.Delete(&models.Permission{}, id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete permission"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Permission deleted"})
	}
}
