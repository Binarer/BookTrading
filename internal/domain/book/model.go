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

// Book представляет собой книгу в системе
type Book struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title" gorm:"not null"`
	Author      string       `json:"author" gorm:"not null"`
	Description string       `json:"description"`
	Photos      []string     `json:"photos" gorm:"-"`        // Игнорируем в GORM
	PhotosJSON  string       `json:"-" gorm:"column:photos"` // Храним как JSON
	UserID      uint         `json:"user_id" gorm:"not null"`
	User        *user.User   `json:"user" gorm:"foreignKey:UserID"`
	StateID     uint         `json:"state_id" gorm:"not null"`
	State       *state.State `json:"state" gorm:"foreignKey:StateID"`
	Tags        []*tag.Tag   `json:"tags" gorm:"many2many:book_tags"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// TableName указывает имя таблицы для модели Book
func (Book) TableName() string {
	return "books"
}

// CreateBookDTO представляет данные для создания книги
type CreateBookDTO struct {
	Title       string   `json:"title" binding:"required"`
	Author      string   `json:"author" binding:"required"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	UserID      uint     `json:"user_id" binding:"required"`
	StateID     uint     `json:"state_id" binding:"required"`
	TagIDs      []uint   `json:"tag_ids"`
}

// UpdateBookDTO представляет данные для обновления книги
type UpdateBookDTO struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	StateID     uint     `json:"state_id"`
	TagIDs      []uint   `json:"tag_ids"`
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
		Photos:      dto.Photos,
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
	if len(dto.Photos) > 0 {
		b.Photos = dto.Photos
	}
	if dto.StateID != 0 {
		b.StateID = dto.StateID
	}
}
