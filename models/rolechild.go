package models

type RoleChild struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	RoleID uint   `json:"role_id"`
	Role   Role   `json:"-"` // Relasi ke Role

}
