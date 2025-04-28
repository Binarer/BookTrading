package mysql

import (
	"booktrading/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLConnection создает новое подключение к MySQL
func NewMySQLConnection(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
} 