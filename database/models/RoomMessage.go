package models

import "gorm.io/gorm"

type RoomMessage struct {
	gorm.Model

	RoomID  string
	UserID  string
	Message string
}
