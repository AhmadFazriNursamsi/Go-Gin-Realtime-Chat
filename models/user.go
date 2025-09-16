package models

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	RoleChild   RoleChild `json:"rolechild"` // Relasi ke RoleChild
	RoleChildID uint      `json:"rolechild_id"`
	Photo       string    `json:"photo"` // 🆕 untuk simpan path file

}
