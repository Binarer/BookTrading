package book

import (
	"booktrading/internal/pkg/gorm"
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/domain/user"
)

// BookState represents the possible states of a book
type BookState string

const (
	StateAvailable BookState = "available"
	StateTrading   BookState = "trading"
	StateTraded    BookState = "traded"
)

// Book represents a book entity
type Book struct {
	gorm.Base
	Title       string `gorm:"size:255;not null"`
	Author      string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Price       float64
	StateID     uint
	State       state.State
	Tags        []tag.Tag `gorm:"many2many:book_tags;"`
	UserID      uint
	User        user.User
}

// TableName specifies the table name for the Book model
func (Book) TableName() string {
	return "books"
}

// CreateBookDTO represents the data needed to create a new book
type CreateBookDTO struct {
	Title       string   `json:"title" validate:"required,min=3,max=100"`
	Author      string   `json:"author" validate:"required,min=3,max=100"`
	Description string   `json:"description" validate:"required,min=10,max=1000"`
	StateID     int64    `json:"state_id" validate:"required"`
	Photos      []string `json:"photos" validate:"required,min=1,max=5,dive,base64"`
	TagIDs      []int64  `json:"tag_ids" validate:"required,min=1,max=5,dive,min=1"`
}

// UpdateBookDTO represents the data needed to update a book
type UpdateBookDTO struct {
	Title       string   `json:"title" validate:"omitempty,min=3,max=100"`
	Author      string   `json:"author" validate:"omitempty,min=3,max=100"`
	Description string   `json:"description" validate:"omitempty,min=10,max=1000"`
	StateID     int64    `json:"state_id" validate:"omitempty,min=1"`
	Photos      []string `json:"photos" validate:"omitempty,min=1,max=5,dive,base64"`
}

// UpdateBookStateDTO represents the data needed to update the state of a book
type UpdateBookStateDTO struct {
	StateID int64 `json:"state_id" validate:"required,min=1"`
} 