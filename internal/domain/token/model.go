package token

import (
	"booktrading/internal/pkg/gorm"
	"time"
)

// RefreshToken представляет собой модель токена обновления
type RefreshToken struct {
	gorm.Base
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName указывает имя таблицы для модели RefreshToken
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
