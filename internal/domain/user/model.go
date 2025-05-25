package user

import (
	"booktrading/internal/pkg/gorm"
	"time"
)

// User представляет собой пользователя системы
// @Description Модель пользователя системы обмена книгами
type User struct {
	gorm.Base
	// @Description ID пользователя
	// @example 1
	ID uint `json:"id" gorm:"primaryKey"`
	// @Description Логин пользователя
	// @example john_doe
	Login string `json:"login" gorm:"type:varchar(50);uniqueIndex;not null"`
	// @Description Отображаемое имя пользователя
	// @example John Doe
	Username string `json:"username" gorm:"type:varchar(50);not null"`
	// @Description Пароль пользователя (не отображается в JSON)
	Password string `json:"-" gorm:"type:varchar(255);not null"`
	// @Description Аватар пользователя в формате base64
	// @example data:image/jpeg;base64,/9j/4AAQSkZJRg...
	Avatar string `json:"avatar,omitempty" gorm:"type:text"`
	// @Description Список ID книг пользователя (не сохраняется в БД)
	BookIDs []uint `json:"book_ids" gorm:"-"`
	// @Description Дата создания
	// @example 2024-03-20T10:00:00Z
	CreatedAt time.Time `json:"created_at"`
	// @Description Дата обновления
	// @example 2024-03-20T10:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName указывает имя таблицы для модели User
func (User) TableName() string {
	return "users"
}

// CreateUserDTO представляет данные для создания пользователя
// @Description Данные для создания нового пользователя
type CreateUserDTO struct {
	// @Description Логин пользователя
	// @example john_doe
	Login string `json:"login" binding:"required,min=3,max=50"`
	// @Description Отображаемое имя пользователя
	// @example John Doe
	Username string `json:"username" binding:"required,min=2,max=50"`
	// @Description Пароль пользователя
	// @example password123
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// UpdateUserDTO представляет данные для обновления пользователя
// @Description Данные для обновления существующего пользователя
type UpdateUserDTO struct {
	// @Description Отображаемое имя пользователя
	// @example John Doe
	Username string `json:"username" binding:"omitempty,min=2,max=50"`
	// @Description Аватар пользователя в формате base64
	// @example data:image/jpeg;base64,/9j/4AAQSkZJRg...
	Avatar string `json:"avatar,omitempty" binding:"omitempty"`
}

// LoginDTO представляет данные для входа пользователя
// @Description Данные для входа в систему
type LoginDTO struct {
	// @Description Логин пользователя
	// @example john_doe
	Login string `json:"login" binding:"required"`
	// @Description Пароль пользователя
	// @example password123
	Password string `json:"password" binding:"required"`
}

// TokenResponse представляет ответ с JWT токеном
// @Description Ответ с токенами доступа
type TokenResponse struct {
	// @Description JWT токен доступа
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
	// @Description Токен для обновления доступа
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refresh_token"`
}

// LoginResponse представляет полный ответ при входе пользователя
// @Description Полный ответ при успешном входе в систему
type LoginResponse struct {
	// @Description JWT токен доступа
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
	// @Description Токен для обновления доступа
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refresh_token"`
	// @Description ID пользователя
	// @example 1
	UserID uint `json:"user_id"`
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
	if dto.Username != "" {
		u.Username = dto.Username
	}
	if dto.Avatar != "" {
		u.Avatar = dto.Avatar
	}
}
