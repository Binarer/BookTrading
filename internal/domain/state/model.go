package state

import (
	"booktrading/internal/pkg/gorm"
)

// State представляет состояние книги
// @Description Модель состояния книги
type State struct {
	gorm.Base
	// @Description Название состояния
	// @example available
	Name string `gorm:"size:50;not null;unique" json:"name"`
}

// TableName указывает имя таблицы для модели State
func (State) TableName() string {
	return "states"
}

// CreateStateDTO представляет данные, необходимые для создания нового состояния
// @Description Данные для создания нового состояния
type CreateStateDTO struct {
	// @Description Название состояния
	// @example available
	Name string `json:"name" validate:"required,min=3,max=50"`
}

// UpdateStateDTO представляет данные, необходимые для обновления состояния
// @Description Данные для обновления существующего состояния
type UpdateStateDTO struct {
	// @Description Название состояния
	// @example available
	Name string `json:"name" validate:"required,min=3,max=50"`
} 