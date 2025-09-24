package models

type Role struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Name      string      `json:"name"`
	RoleChild []RoleChild `json:"rolechild"`
}
