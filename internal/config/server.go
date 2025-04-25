package config

import (
	"strconv"
)

// ServerConfig содержит конфигурацию сервера
type ServerConfig struct {
	Port int
	Host string
}

// NewServerConfig создает новую конфигурацию сервера
func NewServerConfig() *ServerConfig {
	port, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	return &ServerConfig{
		Port: port,
		Host: getEnv("SERVER_HOST", "0.0.0.0"),
	}
} 