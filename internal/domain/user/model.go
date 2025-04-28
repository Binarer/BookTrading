package user

import (
	"booktrading/internal/pkg/gorm"
	"time"
)

// User представляет пользователя системы
// @Description Модель пользователя системы обмена книгами
type User struct {
	gorm.Base
	// @Description Имя пользователя
	// @example john_doe
	Username     string    `gorm:"size:255;not null;unique" json:"username"`
	
	// @Description Email пользователя
	// @example john@example.com
	Email        string    `gorm:"size:255;not null;unique" json:"email"`
	
	// @Description Хеш пароля (не отображается в JSON)
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	
	// @Description Время создания аккаунта
	// @example 2025-04-28T12:00:00Z
	CreatedAt    time.Time `json:"created_at"`
	
	// @Description Время последнего обновления аккаунта
	// @example 2025-04-28T12:00:00Z
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName указывает имя таблицы для модели User
func (User) TableName() string {
	return "users"
} 