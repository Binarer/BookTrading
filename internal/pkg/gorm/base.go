package gorm

import (
	"time"
	"gorm.io/gorm"
)

// Base contains common columns for all tables
type Base struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
} 