package state

import (
	"time"
)

// State represents a book state
type State struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateStateDTO represents the data needed to create a new state
type CreateStateDTO struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

// UpdateStateDTO represents the data needed to update a state
type UpdateStateDTO struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
} 