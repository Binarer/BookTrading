package tag

import (
	"booktrading/internal/pkg/gorm"
)

// Tag represents a book tag
// @Description Tag model for categorizing books
type Tag struct {
	gorm.Base
	// @Description Name of the tag
	// @example fiction
	Name string `gorm:"size:255;not null;unique" json:"name"`
}

// TableName specifies the table name for the Tag model
func (Tag) TableName() string {
	return "tags"
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