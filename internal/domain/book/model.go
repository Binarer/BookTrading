package book

import (
	"booktrading/internal/domain/tag"
	"time"
)

type Book struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Description string     `json:"description"`
	UserID      int64      `json:"user_id"`
	Tags        []*tag.Tag `json:"tags,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
} 