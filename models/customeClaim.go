package models

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Roleid        *uint  `json:"role_id"`
	RoleName      string `json:"role_name"`
	RoleChildID   *uint  `json:"rolechild_id"`
	RoleChildName string `json:"roleChild_name"`
	// RoomsId       uint   `json:"room_id`
	RoomsId []uint // ✅ ubah jadi slice

	jwt.RegisteredClaims
}
