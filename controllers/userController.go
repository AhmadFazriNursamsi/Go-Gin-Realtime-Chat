package controllers

import (
	"net/http"
	// "strconv"
	"myapi/database"
	"myapi/models"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "gorm.io/gorm"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Preload("RoleChild").Find(&users)
	database.DB.Preload("RoleChild.Role").Find(&users)
	c.JSON(http.StatusOK, users)

}
func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Validasi jika ada RoleChildID di input
	if user.RoleChildID != 0 {
		var rc models.RoleChild
		if err := database.DB.First(&rc, user.RoleChildID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RoleChildID"})
			return
		}
	}

	// Simpan ke database
	database.DB.Create(&user)
	// Preload RoleChild supaya response ada detail RoleChild
	database.DB.Preload("RoleChild").Preload("RoleChild.Role").First(&user, user.ID)
	c.JSON(http.StatusCreated, user)

}
func UploadUserPhoto(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Cari user yang mau diupdate
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	// Terima file dari form-data
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo is required"})
		return
	}
	// Buat folder jika belum ada
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	// Simpan file ke folder uploads dengan nama unik
	// filename := filepath.Base(file.Filename)
	filename := uuid.New().String() + filepath.Ext(file.Filename)

	mimeType := file.Header.Get("Content-Type")
	if mimeType != "image/png" && mimeType != "image/jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File harus JPG/PNG"})
		return
	}

	filepath := "uploads/" + filename
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo"})
		return
	}
	// Update path foto di database
	user.Photo = filepath
	database.DB.Save(&user)
	// Preload RoleChild biar response lengkap
	database.DB.Preload("RoleChild").Preload("RoleChild.Role").First(&user, user.ID)
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Cari user yang mau diupdate
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	var input models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Validasi jika ada RoleChildID di input
	if input.RoleChildID != 0 {
		var rc models.RoleChild
		if err := database.DB.First(&rc, input.RoleChildID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RoleChildID"})
			return
		}
	}

	// Update hanya field yang dikirim
	database.DB.Model(&user).Updates(input)

	// Preload RoleChild biar response lengkap
	database.DB.Preload("RoleChild").First(&user, user.ID)

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	database.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
