package mysql

import (
	"booktrading/internal/config"
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/domain/token"
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/logger"
	"fmt"

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
		&book.BookPhoto{},
		&tag.Tag{},
		&state.State{},
		&token.RefreshToken{},
	)
	if err != nil {
		logger.Error("Failed to migrate database", err)
		return nil, err
	}

	logger.Info("Database migration completed successfully")
	return db, nil
}

// InitGormDB создает новое подключение к MySQL через GORM
func InitGormDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
