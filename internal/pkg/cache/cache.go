package cache

import (
	"strings"
	"sync"
	"time"
)

// Cache представляет собой простой in-memory кеш
type Cache struct {
	mu    sync.RWMutex
	items map[string]item
	stop  chan struct{}
}

type item struct {
	value      interface{}
	expiration time.Time
}

// NewCache создает новый экземпляр кеша
func NewCache() *Cache {
	c := &Cache{
		items: make(map[string]item),
		stop:  make(chan struct{}),
	}

	// Запускаем горутину для очистки устаревших записей
	go c.startCleanup()

	return c
}

// startCleanup запускает периодическую очистку устаревших записей
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stop:
			return
		}
	}
}

// cleanup удаляет все устаревшие записи
func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, v := range c.items {
		if now.After(v.expiration) {
			delete(c.items, k)
		}
	}
}

// Set сохраняет значение в кеше
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:      value,
		expiration: time.Now().Add(duration),
	}
}

// Get получает значение из кеша
func (c *Cache) Get(key string) (interface{}, bool) {
	if c == nil {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

// Delete удаляет значение из кеша
func (c *Cache) Delete(key string) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// DeletePattern удаляет все значения, соответствующие паттерну
func (c *Cache) DeletePattern(pattern string) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.items {
		if strings.Contains(k, pattern) {
			delete(c.items, k)
		}
	}
}

// Flush очищает весь кеш
func (c *Cache) Flush() {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]item)
}

// ItemCount возвращает количество элементов в кеше
func (c *Cache) ItemCount() int {
	if c == nil {
		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// Stop останавливает очистку кеша
func (c *Cache) Stop() {
	if c == nil {
		return
	}
	close(c.stop)
}
