package tag

import (
	"booktrading/internal/pkg/gorm"
)

// Tag представляет тег книги
// @Description Модель тега для категоризации книг
type Tag struct {
	gorm.Base
	// @Description Название тега
	// @example fiction
	Name string `gorm:"size:255;not null;unique" json:"name"`
}

// TableName указывает имя таблицы для модели Tag
func (Tag) TableName() string {
	return "tags"
}

// CreateTagDTO представляет данные, необходимые для создания нового тега
// @Description Данные для создания нового тега
type CreateTagDTO struct {
	// @Description Название тега
	// @example fiction
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// UpdateTagDTO представляет данные, необходимые для обновления тега
// @Description Данные для обновления существующего тега
type UpdateTagDTO struct {
	// @Description Название тега
	// @example fiction
	Name string `json:"name" validate:"omitempty,min=1,max=255"`
} 