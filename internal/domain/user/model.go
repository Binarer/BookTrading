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
	Username string `gorm:"size:255;not null;unique" json:"username"`

	// @Description Логин пользователя
	// @example john_doe
	Login string `gorm:"size:255;not null;unique" json:"login"`

	// @Description Хеш пароля (не отображается в JSON)
	PasswordHash string `gorm:"size:255;not null" json:"-"`

	// @Description Время создания аккаунта
	// @example 2025-04-28T12:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// @Description Время последнего обновления аккаунта
	// @example 2025-04-28T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName указывает имя таблицы для модели User
func (User) TableName() string {
	return "users"
}

// CreateUserDTO представляет данные для создания нового пользователя
type CreateUserDTO struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Login    string `json:"login" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserDTO представляет данные для обновления пользователя
type UpdateUserDTO struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
	Login    string `json:"login" validate:"omitempty,min=3,max=50"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

// LoginDTO представляет данные для входа
type LoginDTO struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse представляет ответ с JWT токеном
type TokenResponse struct {
	Token string `json:"token"`
}
