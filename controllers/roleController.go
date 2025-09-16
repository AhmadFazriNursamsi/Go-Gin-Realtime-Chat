package controllers

import (
	"myapi/database"
	"myapi/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {
	var roles []models.Role
	database.DB.Preload("RoleChild").Find(&roles)
	c.JSON(http.StatusOK, roles)
}

func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.BindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&role)
	c.JSON(http.StatusCreated, role)
}

// func AssignRole(c *gin.Context) {
// 	userID := c.Param("userID")
// 	roleID := c.Param("roleID")
// 	var user models.User
// 	if err := database.DB.First(&user, userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
// 		return
// 	}
// 	var role models.Role
// 	if err := database.DB.First(&role, roleID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
// 		return
// 	}
// 	database.DB.Model(&user).Association("Roles").Append(&role)
// 	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to user"})
// }

func DeleteRole(c *gin.Context) {
	userID := c.Param("userID")
	roleID := c.Param("roleID")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	var role models.Role

	if err := database.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
		return
	}
	database.DB.Model(&user).Association("Roles").Delete(&role)
	c.JSON(http.StatusOK, gin.H{"message": "Role removed from user"})
}
func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
		return
	}
	var input models.Role
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Model(&role).Updates(input)
	c.JSON(http.StatusOK, role)
}
