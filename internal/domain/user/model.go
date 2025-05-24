package user

import (
	"booktrading/internal/pkg/gorm"
	"time"
)

// User представляет собой пользователя системы
// @Description Модель пользователя системы обмена книгами
type User struct {
	gorm.Base
	ID        uint      `json:"id" gorm:"primaryKey"`
	Login     string    `json:"login" gorm:"uniqueIndex;not null"`
	Username  string    `json:"username" gorm:"not null"`          // Отображаемое имя пользователя
	Password  string    `json:"-" gorm:"not null"`                 // Не отображаем в JSON
	Avatar    string    `json:"avatar,omitempty" gorm:"type:text"` // Base64 строка для аватарки
	BookIDs   []uint    `json:"book_ids" gorm:"-"`                 // Игнорируем в GORM
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName указывает имя таблицы для модели User
func (User) TableName() string {
	return "users"
}

// CreateUserDTO представляет данные для создания пользователя
type CreateUserDTO struct {
	Login    string `json:"login" binding:"required,min=3,max=50"`
	Username string `json:"username" binding:"required,min=2,max=50"` // Отображаемое имя пользователя
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// UpdateUserDTO представляет данные для обновления пользователя
type UpdateUserDTO struct {
	Login    string `json:"login" binding:"omitempty,min=3,max=50"`
	Username string `json:"username" binding:"omitempty,min=2,max=50"` // Отображаемое имя пользователя
	Avatar   string `json:"avatar,omitempty" binding:"omitempty"`      // Base64 строка для аватарки
}

// LoginDTO представляет данные для входа пользователя
type LoginDTO struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse представляет ответ с JWT токеном
type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// ToUser преобразует DTO в модель User
func (dto *CreateUserDTO) ToUser() *User {
	return &User{
		Login:    dto.Login,
		Username: dto.Username,
		Password: dto.Password,
	}
}

// UpdateFromDTO обновляет поля пользователя из DTO
func (u *User) UpdateFromDTO(dto *UpdateUserDTO) {
	if dto.Login != "" {
		u.Login = dto.Login
	}
	if dto.Username != "" {
		u.Username = dto.Username
	}
	if dto.Avatar != "" {
		u.Avatar = dto.Avatar
	}
}
