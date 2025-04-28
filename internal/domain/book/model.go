package book

import (
	"booktrading/internal/domain/tag"
	"time"
)

// BookState represents the possible states of a book
type BookState string

const (
	StateAvailable BookState = "available"
	StateTrading   BookState = "trading"
	StateTraded    BookState = "traded"
)

// Book represents a book in the system
// @Description Book model for the book trading system
type Book struct {
	// @Description Unique identifier for the book
	// @example 1
	ID          int64      `json:"id"`

	// @Description Title of the book
	// @example The Great Gatsby
	Title       string     `json:"title"`

	// @Description Author of the book
	// @example F. Scott Fitzgerald
	Author      string     `json:"author"`

	// @Description Description of the book
	// @example A story of the fabulously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan.
	Description string     `json:"description"`

	// @Description Current state of the book
	// @example available
	StateID     int64      `json:"state_id"`

	// @Description List of base64 encoded photos of the book
	// @example ["data:image/jpeg;base64,/9j/4AAQSkZJRg..."]
	Photos      []string   `json:"photos"`

	// @Description When the book was added to the system
	// @example 2025-04-28T12:00:00Z
	CreatedAt   time.Time  `json:"created_at"`

	// @Description When the book was last updated
	// @example 2025-04-28T12:00:00Z
	UpdatedAt   time.Time  `json:"updated_at"`

	// @Description List of tags associated with the book
	Tags        []*tag.Tag `json:"tags,omitempty"`
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