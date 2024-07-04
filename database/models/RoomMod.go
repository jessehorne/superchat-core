package models

import "gorm.io/gorm"

type RoomMod struct {
	gorm.Model

	UserID string
	RoomID string
	Role   int
}
