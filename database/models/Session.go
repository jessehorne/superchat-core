package models

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model

	Token     string
	ExpiresAt time.Time
	UserID    string
}
