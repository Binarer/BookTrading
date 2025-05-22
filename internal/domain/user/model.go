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

	// @Description Описание пользователя
	// @example Book lover and collector
	Description *string `gorm:"type:text" json:"description"`

	// @Description Аватар пользователя в формате base64
	// @example data:image/jpeg;base64,/9j/4AAQSkZJRg...
	Avatar *string `gorm:"type:text" json:"avatar"`

	// @Description Книги пользователя
	// @example [{"id": 1, "title": "The Great Gatsby"}]
	Books interface{} `gorm:"-" json:"books"`

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
	Username    string  `json:"username" validate:"omitempty,min=3,max=50"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	Avatar      *string `json:"avatar" validate:"omitempty,base64"`
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
