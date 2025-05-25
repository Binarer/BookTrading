package book

import (
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/domain/user"
	"time"
)

// BookState представляет возможные состояния книги
// @Description Перечисление возможных состояний книги
type BookState string

const (
	StateAvailable BookState = "available" // Книга доступна для обмена
	StateTrading   BookState = "trading"   // Книга в процессе обмена
	StateTraded    BookState = "traded"    // Книга обменена
)

// BookPhoto представляет фотографию книги
// @Description Модель фотографии книги
type BookPhoto struct {
	// @Description ID фотографии
	// @example 1
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned"`
	// @Description ID книги, к которой относится фотография
	// @example 1
	BookID uint `json:"book_id" gorm:"not null;index;type:int unsigned"`
	// @Description URL фотографии в формате base64
	// @example data:image/jpeg;base64,/9j/4AAQSkZJRg...
	PhotoURL string `json:"photo_url" gorm:"type:mediumtext;not null"`
	// @Description Флаг, указывающий является ли фотография главной
	// @example true
	IsMain bool `json:"is_main" gorm:"default:false"`
	// @Description Дата создания
	// @example 2024-03-20T10:00:00Z
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	// @Description Дата обновления
	// @example 2024-03-20T10:00:00Z
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName указывает имя таблицы для модели BookPhoto
func (BookPhoto) TableName() string {
	return "book_photos"
}

// Book представляет собой книгу в системе
// @Description Модель книги в системе обмена
type Book struct {
	// @Description ID книги
	// @example 1
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned"`
	// @Description Название книги
	// @example Война и мир
	Title string `json:"title" gorm:"type:varchar(255);not null;index"`
	// @Description Автор книги
	// @example Лев Толстой
	Author string `json:"author" gorm:"type:varchar(255);not null;index"`
	// @Description Описание книги
	// @example Роман-эпопея, описывающий русское общество в эпоху войн против Наполеона
	Description string `json:"description" gorm:"type:mediumtext"`
	// @Description ID владельца книги
	// @example 1
	UserID uint `json:"user_id" gorm:"not null;type:int unsigned;index"`
	// @Description Информация о владельце книги
	User *user.User `json:"user" gorm:"foreignKey:UserID"`
	// @Description ID состояния книги
	// @example 1
	StateID uint `json:"state_id" gorm:"not null;type:int unsigned;index"`
	// @Description Информация о состоянии книги
	State *state.State `json:"state" gorm:"foreignKey:StateID"`
	// @Description Теги книги
	Tags []*tag.Tag `json:"tags" gorm:"many2many:book_tags"`
	// @Description Фотографии книги
	Photos []*BookPhoto `json:"photos" gorm:"foreignKey:BookID"`
	// @Description Дата создания
	// @example 2024-03-20T10:00:00Z
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	// @Description Дата обновления
	// @example 2024-03-20T10:00:00Z
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName указывает имя таблицы для модели Book
func (Book) TableName() string {
	return "books"
}

// BookPhotoData представляет данные фотографии для создания книги
// @Description Данные фотографии для создания книги
type BookPhotoData struct {
	// @Description URL фотографии в формате base64
	// @example "data:image/jpeg;base64,/9j/4AAQSkZJRg..."
	PhotoURL string `json:"photo_url"`

	// @Description Флаг, указывающий является ли фотография главной
	// @example true
	IsMain bool `json:"is_main"`
}

// CreateBookDTO представляет данные для создания книги
// @Description Данные для создания новой книги
type CreateBookDTO struct {
	// @Description Название книги
	// @example "Война и мир"
	Title string `json:"title" binding:"required"`

	// @Description Автор книги
	// @example "Лев Толстой"
	Author string `json:"author" binding:"required"`

	// @Description Описание книги
	// @example "Роман-эпопея, описывающий русское общество в эпоху войн против Наполеона"
	Description string `json:"description"`

	// @Description Массив фотографий книги
	Photos []BookPhotoData `json:"photos"`

	// @Description ID пользователя-владельца книги
	// @example 1
	UserID uint `json:"user_id" binding:"required"`

	// @Description ID состояния книги (1 - доступна, 2 - недоступна)
	// @example 1
	StateID uint `json:"state_id" binding:"required"`

	// @Description Массив ID тегов книги
	// @example [1, 2, 3]
	TagIDs []uint `json:"tag_ids"`
}

// UpdateBookDTO представляет данные для обновления книги
// @Description Данные для обновления существующей книги
type UpdateBookDTO struct {
	// @Description Название книги
	// @example Война и мир
	Title string `json:"title"`
	// @Description Автор книги
	// @example Лев Толстой
	Author string `json:"author"`
	// @Description Описание книги
	// @example Роман-эпопея, описывающий русское общество в эпоху войн против Наполеона
	Description string `json:"description"`
	// @Description Массив URL фотографий в формате base64
	// @example ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."]
	Photos []string `json:"photos"`
	// @Description ID состояния книги
	// @example 1
	StateID uint `json:"state_id"`
	// @Description Массив ID тегов книги
	// @example [1, 2, 3]
	TagIDs []uint `json:"tag_ids"`
}

// UpdateBookStateDTO представляет данные, необходимые для обновления состояния книги
// @Description Данные для обновления состояния книги
type UpdateBookStateDTO struct {
	// @Description ID нового состояния книги
	// @example 1
	StateID int64 `json:"state_id" validate:"required,min=1"`
}

// ToBook преобразует DTO в модель Book
func (dto *CreateBookDTO) ToBook() *Book {
	return &Book{
		Title:       dto.Title,
		Author:      dto.Author,
		Description: dto.Description,
		UserID:      dto.UserID,
		StateID:     dto.StateID,
	}
}

// UpdateFromDTO обновляет поля книги из DTO
func (b *Book) UpdateFromDTO(dto *UpdateBookDTO) {
	if dto.Title != "" {
		b.Title = dto.Title
	}
	if dto.Author != "" {
		b.Author = dto.Author
	}
	if dto.Description != "" {
		b.Description = dto.Description
	}
	if dto.StateID != 0 {
		b.StateID = dto.StateID
	}
}
