package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache представляет собой обертку над go-cache
type Cache struct {
	cache *cache.Cache
}

// NewCache создает новый экземпляр кеша
func NewCache(defaultTTL, cleanupInterval time.Duration) *Cache {
	return &Cache{
		cache: cache.New(defaultTTL, cleanupInterval),
	}
}

// Set сохраняет значение в кеше
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.cache.Set(key, value, ttl)
}

// Get получает значение из кеша
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// Delete удаляет значение из кеша
func (c *Cache) Delete(key string) {
	c.cache.Delete(key)
}

// Flush очищает весь кеш
func (c *Cache) Flush() {
	c.cache.Flush()
}

// ItemCount возвращает количество элементов в кеше
func (c *Cache) ItemCount() int {
	return c.cache.ItemCount()
} 