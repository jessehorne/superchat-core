package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model

	Name              string
	PasswordProtected bool
	Password          string
	PasswordSalt      string
}
