package tag

import "time"

// Tag represents a book tag
// @Description Tag model for categorizing books
type Tag struct {
	// @Description Unique identifier for the tag
	// @example 1
	ID        int64     `json:"id"`

	// @Description Name of the tag
	// @example fiction
	Name      string    `json:"name"`

	// @Description When the tag was created
	// @example 2025-04-28T12:00:00Z
	CreatedAt time.Time `json:"created_at"`

	// @Description When the tag was last updated
	// @example 2025-04-28T12:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTagDTO represents the data needed to create a new tag
type CreateTagDTO struct {
	// @Description Name of the tag
	// @example fiction
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// UpdateTagDTO represents the data needed to update a tag
type UpdateTagDTO struct {
	// @Description Name of the tag
	// @example fiction
	Name string `json:"name" validate:"omitempty,min=1,max=255"`
} 