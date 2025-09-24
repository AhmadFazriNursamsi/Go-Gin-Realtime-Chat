package models

type RolePermission struct {
	ID           uint  `gorm:"primaryKey" json:"id"`
	RoleID       *uint `json:"role_id"`
	RoleChildID  *uint `json:"rolechild_id"` // boleh null
	PermissionID uint  `json:"permission_id"`

	RoleChild  RoleChild  `gorm:"foreignKey:RoleChildID"`
	Role       Role       `gorm:"foreignKey:RoleID"`
	Permission Permission `gorm:"foreignKey:PermissionID"`
}
