package gorm

import (
	"time"
	"gorm.io/gorm"
)

// Base содержит общие поля для всех таблиц
// @Description Базовая структура для всех моделей
type Base struct {
	// @Description Уникальный идентификатор
	// @example 1
	ID        uint           `gorm:"primarykey" json:"id"`
	
	// @Description Время создания записи
	// @example 2025-04-28T12:00:00Z
	CreatedAt time.Time      `json:"created_at"`
	
	// @Description Время последнего обновления записи
	// @example 2025-04-28T12:00:00Z
	UpdatedAt time.Time      `json:"updated_at"`
	
	// @Description Время удаления записи (для мягкого удаления)
	// @example 2025-04-28T12:00:00Z
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
} 