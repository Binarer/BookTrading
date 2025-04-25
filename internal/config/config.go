package config

import (
	"os"
	_ "path/filepath"
	"strconv"
	_ "time"

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
	}, nil
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
