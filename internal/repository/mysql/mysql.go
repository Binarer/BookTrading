package mysql

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB инициализирует подключение к базе данных и выполняет миграции
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", err)
		return nil, err
	}

	// Автоматическая миграция моделей
	err = db.AutoMigrate(
		&user.User{},
		&book.Book{},
		&tag.Tag{},
		&state.State{},
	)
	if err != nil {
		logger.Error("Failed to migrate database", err)
		return nil, err
	}

	logger.Info("Database migration completed successfully")
	return db, nil
}
