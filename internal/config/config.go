package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/rs/zerolog"
)

// Config содержит все конфигурации приложения
type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	Cache    *CacheConfig
	Logging  *LoggingConfig
	CORS     *CORSConfig
	JWT      *JWTConfig
}

// ServerConfig содержит конфигурацию сервера
type ServerConfig struct {
	Address string
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

// NewConfig создает новую конфигурацию приложения
func NewConfig() (*Config, error) {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cacheConfig, err := NewCacheConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Server:   NewServerConfig(),
		Database: NewDatabaseConfig(),
		Cache:    cacheConfig,
		Logging:  NewLoggingConfig(),
		CORS:     NewCORSConfig(),
		JWT:      &JWTConfig{SecretKey: getEnv("JWT_SECRET_KEY", "your-secret-key")},
	}, nil
}

// NewServerConfig создает новую конфигурацию сервера
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Address: getEnv("SERVER_ADDRESS", ":8000"),
	}
}

// NewDatabaseConfig creates a new database configuration from environment variables
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

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt получает целочисленное значение переменной окружения или возвращает значение по умолчанию
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
