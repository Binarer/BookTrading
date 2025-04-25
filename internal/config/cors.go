package config

import "strconv"

// CORSConfig содержит конфигурацию CORS
type CORSConfig struct {
	AllowedOrigins   string
	AllowedMethods   string
	AllowedHeaders   string
	ExposedHeaders   string
	AllowCredentials bool
	MaxAge           int
}

// NewCORSConfig создает новую конфигурацию CORS
func NewCORSConfig() *CORSConfig {
	maxAge, _ := strconv.Atoi(getEnv("CORS_MAX_AGE", "300"))
	allowCredentials, _ := strconv.ParseBool(getEnv("CORS_ALLOW_CREDENTIALS", "true"))

	return &CORSConfig{
		AllowedOrigins:   getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:8000,http://10.3.13.28:8000"),
		AllowedMethods:   getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		AllowedHeaders:   getEnv("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token"),
		ExposedHeaders:   getEnv("CORS_EXPOSED_HEADERS", "Link"),
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
	}
}
