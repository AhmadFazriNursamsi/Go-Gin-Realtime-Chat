package models

import "time"

type Messages struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoomID     uint      `json:"room_id"`
	SenderID   uint      `json:"sender_id"`
	Content    string    `json:"content"`
	Type       string    `json:"type"` // text/image/...
	CreatedAt  time.Time `json:"created_at"`
	SenderName string    `json:"sender_name" gorm:"-"`
}
