package models

type Role struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Name      string      `json:"name"`
	GUIDRole  string      `json:"guid_role"`
	RoleChild []RoleChild `json:"rolechild"`
}
