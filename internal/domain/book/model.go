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
type BookPhoto struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned"`
	BookID    uint      `json:"book_id" gorm:"not null;index;type:int unsigned"`
	PhotoURL  string    `json:"photo_url" gorm:"type:mediumtext;not null"`
	IsMain    bool      `json:"is_main" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName указывает имя таблицы для модели BookPhoto
func (BookPhoto) TableName() string {
	return "book_photos"
}

// Book представляет собой книгу в системе
type Book struct {
	ID          uint         `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned"`
	Title       string       `json:"title" gorm:"type:varchar(255);not null;index"`
	Author      string       `json:"author" gorm:"type:varchar(255);not null;index"`
	Description string       `json:"description" gorm:"type:mediumtext"`
	Photos      []BookPhoto  `json:"photos" gorm:"foreignKey:BookID"`
	UserID      uint         `json:"user_id" gorm:"not null;type:int unsigned;index"`
	User        *user.User   `json:"user" gorm:"foreignKey:UserID"`
	StateID     uint         `json:"state_id" gorm:"not null;type:int unsigned;index"`
	State       *state.State `json:"state" gorm:"foreignKey:StateID"`
	Tags        []*tag.Tag   `json:"tags" gorm:"many2many:book_tags"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
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
	book := &Book{
		Title:       dto.Title,
		Author:      dto.Author,
		Description: dto.Description,
		UserID:      dto.UserID,
		StateID:     dto.StateID,
	}

	// Создаем фотографии
	if len(dto.Photos) > 0 {
		book.Photos = make([]BookPhoto, len(dto.Photos))
		for i, photoURL := range dto.Photos {
			book.Photos[i] = BookPhoto{
				PhotoURL: photoURL,
				IsMain:   i == 0, // Первая фотография - главная
			}
		}
	}

	return book
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
		b.Photos = make([]BookPhoto, len(dto.Photos))
		for i, photoURL := range dto.Photos {
			b.Photos[i] = BookPhoto{
				PhotoURL: photoURL,
				IsMain:   i == 0, // Первая фотография - главная
			}
		}
	}
	if dto.StateID != 0 {
		b.StateID = dto.StateID
	}
}
