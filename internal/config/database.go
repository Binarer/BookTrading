package config

import (
	"strconv"
)

// DatabaseConfig содержит конфигурацию базы данных
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// NewDatabaseConfig создает новую конфигурацию базы данных
func NewDatabaseConfig() *DatabaseConfig {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "mysql"),
		Port:     port,
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "booktrading"),
	}
}
