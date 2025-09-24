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
	"gorm.io/gorm"
	// "gorm.io/gorm"
)

type UserWithPermissions struct {
	ID            uint     `json:"id"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	RoleName      string   `json:"role_name"`
	RoleChileName string   `json:"role_child_name"`
	Permissions   []string `json:"permissions"`
}

func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Preload("RoleChild").Find(&users)
	database.DB.Preload("RoleChild.Role").Find(&users)
	c.JSON(http.StatusOK, users)

}

type tokens struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}
type ErrorResponses struct {
	Message string `json:"message" example:"Email atau username tidak ditemukan"`
}

// LoginHandler godoc
// @Summary Login user
// @Description Login untuk mendapatkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login"
// @Success 200 {object} tokens
// @Failure 400 {object} ErrorResponses
// @Router /login [post]
// @Security Bearer
func GetUsersWithPermissions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.Preload("Role").Preload("RoleChild").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
			return
		}

		var result []UserWithPermissions

		for _, u := range users {
			var perms []string

			query := db.Table("role_permissions").
				Select("permissions.name").
				Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
				Where("role_permissions.role_id = ?", u.Role.ID)

			if u.RoleChildID != nil && *u.RoleChildID != 0 {
				// Jika ingin AND (harus cocok dua-duanya)
				query = query.Where("role_permissions.role_child_id = ?", *u.RoleChildID)

				// Jika ingin OR (salah satu cocok boleh)
				// query = query.Or("role_permissions.role_child_id = ?", *u.RoleChildID)
			}

			query.Pluck("permissions.name", &perms)

			// ✅ Hapus duplikat
			permMap := make(map[string]bool)
			uniquePerms := make([]string, 0, len(perms))
			for _, p := range perms {
				if !permMap[p] {
					permMap[p] = true
					uniquePerms = append(uniquePerms, p)
				}
			}

			result = append(result, UserWithPermissions{
				ID:            u.ID,
				Name:          u.Name,
				Email:         u.Email,
				RoleName:      u.Role.Name,
				RoleChileName: u.RoleChild.Name,
				Permissions:   uniquePerms,
			})
		}

		c.JSON(http.StatusOK, result)
	}
}

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Validasi jika ada RoleChildID di input
	if user.RoleChildID != nil {
		var rc models.RoleChild
		if err := database.DB.First(&rc, *user.RoleChildID).Error; err != nil {
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
	if input.RoleChildID != nil {
		var rc models.RoleChild
		if err := database.DB.First(&rc, *input.RoleChildID).Error; err != nil {
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
