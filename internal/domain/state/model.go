package state

import (
	"booktrading/internal/pkg/gorm"
)

// State represents a book state
type State struct {
	gorm.Base
	Name string `gorm:"size:50;not null;unique"`
}

// TableName specifies the table name for the State model
func (State) TableName() string {
	return "states"
}

// CreateStateDTO represents the data needed to create a new state
type CreateStateDTO struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

// UpdateStateDTO represents the data needed to update a state
type UpdateStateDTO struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
} 