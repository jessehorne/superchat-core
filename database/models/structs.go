package models

import (
	"gorm.io/gorm"
	"time"
)

type GivenFields struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
