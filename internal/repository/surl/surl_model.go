package surl

import (
	"gorm.io/gorm"
)

type ShortURL struct {
	gorm.Model
	ShortCode   string `gorm:"uniqueIndex"`
	OriginalURL string `gorm:"not null"`
	UserID      uint
}
