package models

import (
	"time"
)

type Session struct {
	GivenFields

	Token     string
	ExpiresAt time.Time
	UserID    string
}
