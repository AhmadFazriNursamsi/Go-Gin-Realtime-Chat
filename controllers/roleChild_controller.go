package controllers

import (
	"myapi/database"
	"myapi/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Getrolechild(c *gin.Context) {
	var rolechild []models.RoleChild
	database.DB.Preload("Role").Find(&rolechild)
	c.JSON(http.StatusOK, rolechild)

}
func Createrolechild(c *gin.Context) {
	var rc models.RoleChild

	if err := c.BindJSON(&rc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan ke database
	database.DB.Create(&rc)

	// Preload Role supaya response ada detail Role
	database.DB.Preload("Role").First(&rc, rc.ID)

	c.JSON(http.StatusCreated, rc)
}
func Updaterolechild(c *gin.Context) {
	id := c.Param("id")
	var rolechild models.RoleChild
	if err := database.DB.First(&rolechild, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "RoleChild not found"})
		return
	}
	var input models.RoleChild
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Model(&rolechild).Updates(input)
	database.DB.Preload("Role").First(&rolechild, rolechild.ID)

	c.JSON(http.StatusOK, rolechild)
	// Preload Role supaya response ada detail Role
	// database.DB.Preload("Role").First(&rolechild, rolechild.ID)
	// // database.DB.Model(&rolechild).Updates(input)
	// database.DB.Preload("").Updates(input)
	// c.JSON(http.StatusOK, rolechild)
}
func Deleterolechild(c *gin.Context) {
	id := c.Param("id")
	var rolechild models.RoleChild
	if err := database.DB.First(&rolechild, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "RoleChild not found"})
		return
	}
	database.DB.Delete(&rolechild)
	c.JSON(http.StatusOK, gin.H{"message": "RoleChild deleted"})
}
