package main

import (
	"booktrading/internal/config"
	"booktrading/internal/pkg/logger"
	"booktrading/internal/repository"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Инициализация логгера
	logger.Init()

	// Получение конфигурации БД
	dbConfig := config.NewDatabaseConfig()
	db, err := repository.NewMySQLConnection(dbConfig)
	if err != nil {
		logger.Error("Failed to connect to database", err)
		os.Exit(1)
	}
	defer db.Close()

	// Получение списка миграций
	migrations, err := getMigrations()
	if err != nil {
		logger.Error("Failed to get migrations", err)
		os.Exit(1)
	}

	// Создание таблицы для отслеживания миграций
	if err := createMigrationsTable(db); err != nil {
		logger.Error("Failed to create migrations table", err)
		os.Exit(1)
	}

	// Применение миграций
	for _, migration := range migrations {
		if err := applyMigration(db, migration); err != nil {
			logger.Error(fmt.Sprintf("Failed to apply migration %s", migration), err)
			os.Exit(1)
		}
	}

	logger.Info("All migrations applied successfully")
}

func getMigrations() ([]string, error) {
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func applyMigration(db *sql.DB, migration string) error {
	// Проверяем, была ли уже применена миграция
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", migration).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		logger.Info(fmt.Sprintf("Migration %s already applied, skipping", migration))
		return nil
	}

	// Читаем SQL файл
	content, err := ioutil.ReadFile(filepath.Join("migrations", migration))
	if err != nil {
		return err
	}

	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Применяем миграцию
	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	// Записываем информацию о примененной миграции
	if _, err := tx.Exec("INSERT INTO migrations (name) VALUES (?)", migration); err != nil {
		return err
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Successfully applied migration %s", migration))
	return nil
}
