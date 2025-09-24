package controllers

import (
	"myapi/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssignPermissionToRole(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			RoleChild    *int `json:"role_child"`
			RoleID       *int `json:"role_id"`
			PermissionID int  `json:"permission_id"`
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		// Cek apakah sudah ada
		var rp models.RolePermission

		dbQuery := db.Where("permission_id = ?", req.PermissionID)
		if req.RoleChild != nil {
			dbQuery = dbQuery.Where("role_child_id = ?", *req.RoleChild)
		}
		if req.RoleID != nil {
			dbQuery = dbQuery.Where("role_id = ?", *req.RoleID)
		}
		if err := dbQuery.First(&rp).Error; err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Permission already assigned to this role"})
			return
		}

		var roleChildID, roleID uint
		if req.RoleChild != nil {
			roleChildID = uint(*req.RoleChild)
		}
		if req.RoleID != nil {
			roleID = uint(*req.RoleID)
		}
		rp = models.RolePermission{
			RoleChildID:  &roleChildID,
			RoleID:       &roleID,
			PermissionID: uint(req.PermissionID),
		}

		if err := db.Create(&rp).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign permission"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Permission assigned to role"})
	}
}

func RemovePermissionFromRole(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleID, _ := strconv.Atoi(ctx.Param("roleID"))
		permID, _ := strconv.Atoi(ctx.Param("permissionID"))

		if err := db.Where("role_id = ? AND permission_id = ?", roleID, permID).
			Delete(&models.RolePermission{}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove permission"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Permission removed from role"})
	}
}
