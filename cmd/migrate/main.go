package main

import (
	"booktrading/internal/config"
	"booktrading/internal/pkg/logger"
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
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", err)
	}

	// Создание строки подключения к MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
	)

	// Подключение к MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("Failed to connect to MySQL", err)
	}
	defer db.Close()

	// Создание базы данных, если она не существует
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DB.Name))
	if err != nil {
		logger.Fatal("Failed to create database", err)
	}

	// Использование базы данных
	_, err = db.Exec(fmt.Sprintf("USE %s", cfg.DB.Name))
	if err != nil {
		logger.Fatal("Failed to use database", err)
	}

	// Получение списка миграций
	migrations, err := getMigrations()
	if err != nil {
		logger.Fatal("Failed to get migrations", err)
	}

	// Создание таблицы для отслеживания миграций
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		logger.Fatal("Failed to create migrations table", err)
	}

	// Применение миграций
	for _, migration := range migrations {
		// Проверка, была ли миграция уже применена
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE id = ?", migration.Name).Scan(&count)
		if err != nil {
			logger.Fatal("Failed to check migration status", err)
		}

		if count > 0 {
			logger.Info(fmt.Sprintf("Migration %s already applied, skipping", migration.Name))
			continue
		}

		// Чтение SQL файла
		sqlBytes, err := ioutil.ReadFile(migration.Path)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Failed to read migration file %s", migration.Name), err)
		}

		// Выполнение миграции
		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			logger.Fatal(fmt.Sprintf("Failed to apply migration %s", migration.Name), err)
		}

		// Запись информации о примененной миграции
		_, err = db.Exec("INSERT INTO migrations (id) VALUES (?)", migration.Name)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Failed to record migration %s", migration.Name), err)
		}

		logger.Info(fmt.Sprintf("Successfully applied migration %s", migration.Name))
	}

	logger.Info("All migrations completed successfully")
}

type Migration struct {
	Name string
	Path string
}

func getMigrations() ([]Migration, error) {
	var migrations []Migration

	// Получение списка SQL файлов в директории migrations
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, Migration{
				Name: file.Name(),
				Path: filepath.Join("migrations", file.Name()),
			})
		}
	}

	// Сортировка миграций по имени файла
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	return migrations, nil
} 