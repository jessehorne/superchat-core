package models

import "gorm.io/gorm"

type RoomUser struct {
	gorm.Model

	RoomID string
	UserID string
	Muted  bool
}
