package config

import (
	"github.com/rs/zerolog"
)

// LoggingConfig содержит конфигурацию логирования
type LoggingConfig struct {
	Level  zerolog.Level
	Format string
}

// NewLoggingConfig создает новую конфигурацию логирования
func NewLoggingConfig() *LoggingConfig {
	levelStr := getEnv("LOG_LEVEL", "debug")
	level, _ := zerolog.ParseLevel(levelStr)
	
	return &LoggingConfig{
		Level:  level,
		Format: getEnv("LOG_FORMAT", "json"),
	}
} 