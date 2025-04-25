package config

import "time"

// CacheConfig содержит конфигурацию кэша
type CacheConfig struct {
	TTL              time.Duration
	CleanupInterval  time.Duration
}

// NewCacheConfig создает новую конфигурацию кэша
func NewCacheConfig() (*CacheConfig, error) {
	ttl, err := time.ParseDuration(getEnv("CACHE_TTL", "5m"))
	if err != nil {
		return nil, err
	}

	cleanupInterval, err := time.ParseDuration(getEnv("CACHE_CLEANUP_INTERVAL", "10m"))
	if err != nil {
		return nil, err
	}

	return &CacheConfig{
		TTL:             ttl,
		CleanupInterval: cleanupInterval,
	}, nil
} 