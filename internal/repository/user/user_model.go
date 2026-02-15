package user

import (
	"github.com/unify-z/go-surl/internal/repository/surl"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Username     string          `gorm:"uniqueIndex;not null"`
	PasswordHash string          `gorm:"not null"`
	Email        string          `gorm:"uniqueIndex;not null"`
	IsAdmin      bool            `gorm:"default:false"`
	IsBanned     bool            `gorm:"default:false"`
	URLs         []surl.ShortURL `gorm:"foreignKey:UserID"`
}
