package repository

import (
	"booktrading/internal/config"
	"booktrading/internal/domain/repository"
	"booktrading/internal/repository/mysql"
	"database/sql"
	"strconv"

	"gorm.io/gorm"
)

// NewMySQLConnection создает новое подключение к MySQL
func NewMySQLConnection(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + ")/" + cfg.DBName + "?parseTime=true"
	return sql.Open("mysql", dsn)
}

func NewRepository(db *gorm.DB) *repository.Repository {
	return &repository.Repository{
		User:  mysql.NewUserRepository(db),
		Book:  mysql.NewBookRepository(db),
		Tag:   mysql.NewTagRepository(db),
		State: mysql.NewStateRepository(db),
		Token: mysql.NewRefreshTokenRepository(db),
	}
}
