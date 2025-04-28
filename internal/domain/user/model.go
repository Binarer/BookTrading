package user

import (
	"booktrading/internal/pkg/gorm"
	"time"
)

type User struct {
	gorm.Base
	Username     string    `gorm:"size:255;not null;unique" json:"username"`
	Email        string    `gorm:"size:255;not null;unique" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
} 