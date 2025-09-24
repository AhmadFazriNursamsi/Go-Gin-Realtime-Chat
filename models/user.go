package models

type User struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    string  `json:"-"`            // disembunyikan dari JSON
	RoleID      *uint   `json:"role_id"`      // boleh null
	RoleChildID *uint   `json:"rolechild_id"` // boleh null
	Rooms       []Rooms `gorm:"many2many:room_members;"`

	RoleChild RoleChild `gorm:"foreignKey:RoleChildID"`
	Role      Role      `gorm:"foreignKey:RoleID"`
	Photo     string    `json:"photo"` // 🆕 untuk simpan path file
}
