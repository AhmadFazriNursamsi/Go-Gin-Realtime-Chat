package models

import "time"

type Rooms struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	IsGroupChat bool      `json:"is_group_chat"`
	Createdat   time.Time `json:"created_at"`
	Users       []User    `gorm:"many2many:room_members;"`
}
