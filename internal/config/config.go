package config

import (
	"booktrading/internal/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/rs/zerolog"
	"os"
	"strconv"
)

// Config содержит все конфигурации приложения
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    *CacheConfig
	Logging  *LoggingConfig
	CORS     *CORSConfig
	JWT      JWTConfig
}

// ServerConfig содержит конфигурацию сервера
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// JWTConfig represents the JWT configuration
type JWTConfig struct {
	SecretKey string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file", err)
		return nil, err
	}

	// Загрузка конфигурации сервера
	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8000"))
	if err != nil {
		logger.Error("Failed to parse SERVER_PORT", err)
		return nil, err
	}

	// Загрузка конфигурации базы данных
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		logger.Error("Failed to parse DB_PORT", err)
		return nil, err
	}

	cacheConfig, err := NewCacheConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: serverPort,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "booktrading"),
		},
		Cache:   cacheConfig,
		Logging: NewLoggingConfig(),
		CORS:    NewCORSConfig(),
		JWT:     JWTConfig{SecretKey: getEnv("JWT_SECRET_KEY", "your-secret-key-here")},
	}, nil
}

// NewServerConfig создает новую конфигурацию сервера
func NewServerConfig() ServerConfig {
	port, _ := strconv.Atoi(getEnv("SERVER_PORT", "8000"))
	return ServerConfig{
		Host: getEnv("SERVER_HOST", "localhost"),
		Port: port,
	}
}

// NewDatabaseConfig creates a new database configuration from environment variables
func NewDatabaseConfig() DatabaseConfig {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		DBName:   getEnv("DB_NAME", "booktrading"),
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
