package book

import (
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/gorm"
)

// BookState представляет возможные состояния книги
// @Description Перечисление возможных состояний книги
type BookState string

const (
	StateAvailable BookState = "available" // Книга доступна для обмена
	StateTrading   BookState = "trading"   // Книга в процессе обмена
	StateTraded    BookState = "traded"    // Книга обменена
)

// Book представляет сущность книги
// @Description Модель книги для системы обмена книгами
type Book struct {
	gorm.Base
	// @Description Название книги
	// @example The Great Gatsby
	Title string `gorm:"size:255;not null" json:"title"`

	// @Description Автор книги
	// @example F. Scott Fitzgerald
	Author string `gorm:"size:255;not null" json:"author"`

	// @Description Описание книги
	// @example A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.
	Description string `gorm:"type:text" json:"description"`

	// @Description Цена книги
	// @example 19.99
	Price float64 `json:"price"`

	// @Description ID состояния книги
	// @example 1
	StateID uint `json:"state_id"`

	// @Description Состояние книги
	State state.State `json:"state"`

	// @Description Теги книги
	Tags []tag.Tag `gorm:"many2many:book_tags;" json:"tags"`

	// @Description ID пользователя-владельца
	// @example 1
	UserID uint `json:"user_id"`

	// @Description Пользователь-владелец
	User user.User `json:"user"`

	// @Description Фотографии книги в формате base64
	// @example ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."]
	Photos string `gorm:"type:json" json:"photos"`
}

// TableName указывает имя таблицы для модели Book
func (Book) TableName() string {
	return "books"
}

// CreateBookDTO представляет данные, необходимые для создания новой книги
// @Description Данные для создания новой книги
type CreateBookDTO struct {
	// @Description Название книги
	// @example The Great Gatsby
	Title string `json:"title" validate:"required,min=3,max=100"`

	// @Description Автор книги
	// @example F. Scott Fitzgerald
	Author string `json:"author" validate:"required,min=3,max=100"`

	// @Description Описание книги
	// @example A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.
	Description string `json:"description" validate:"required,min=10,max=1000"`

	// @Description ID состояния книги
	// @example 1
	StateID int64 `json:"state_id" validate:"required"`

	// @Description Фотографии книги в формате base64
	// @example ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."]
	Photos []string `json:"photos" validate:"required,min=1,max=5,dive,base64"`

	// @Description ID тегов книги
	// @example [1, 2, 3]
	TagIDs []int64 `json:"tag_ids" validate:"required,min=1,max=5,dive,min=1"`
}

// UpdateBookDTO представляет данные, необходимые для обновления книги
// @Description Данные для обновления существующей книги
type UpdateBookDTO struct {
	// @Description Название книги
	// @example The Great Gatsby
	Title string `json:"title" validate:"omitempty,min=3,max=100"`

	// @Description Автор книги
	// @example F. Scott Fitzgerald
	Author string `json:"author" validate:"omitempty,min=3,max=100"`

	// @Description Описание книги
	// @example A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.
	Description string `json:"description" validate:"omitempty,min=10,max=1000"`

	// @Description ID состояния книги
	// @example 1
	StateID int64 `json:"state_id" validate:"omitempty,min=1"`

	// @Description Фотографии книги в формате base64
	// @example ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."]
	Photos []string `json:"photos" validate:"omitempty,min=1,max=5,dive,base64"`
}

// UpdateBookStateDTO представляет данные, необходимые для обновления состояния книги
// @Description Данные для обновления состояния книги
type UpdateBookStateDTO struct {
	// @Description ID нового состояния книги
	// @example 1
	StateID int64 `json:"state_id" validate:"required,min=1"`
}
