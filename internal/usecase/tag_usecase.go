package usecase

import (
	"booktrading/internal/domain/tag"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/pkg/logger"
	"booktrading/internal/repository"
	_ "encoding/json"
	"fmt"
	"time"
)

// TagUsecase определяет интерфейс для работы с тегами
type TagUsecase interface {
	CreateTag(name string) (*tag.Tag, error)
	GetTagByID(id int64) (*tag.Tag, error)
	GetTagByName(name string) (*tag.Tag, error)
	GetPopularTags(limit int) ([]*tag.Tag, error)
}

// tagUsecase реализует интерфейс TagUsecase
type tagUsecase struct {
	tagRepo repository.TagRepository
	cache   *cache.Cache
}

// NewTagUsecase создает новый экземпляр tagUsecase
func NewTagUsecase(tagRepo repository.TagRepository, cache *cache.Cache) TagUsecase {
	return &tagUsecase{
		tagRepo: tagRepo,
		cache:   cache,
	}
}

// CreateTag создает новый тег
func (u *tagUsecase) CreateTag(name string) (*tag.Tag, error) {
	t := &tag.Tag{Name: name}
	if err := u.tagRepo.Create(t); err != nil {
		logger.Error("Failed to create tag in repository", err)
		return nil, err
	}

	// Инвалидация кеша популярных тегов
	u.cache.Delete("popular_tags")

	return t, nil
}

// GetTagByID получает тег по ID
func (u *tagUsecase) GetTagByID(id int64) (*tag.Tag, error) {
	// Попытка получить тег из кеша
	cacheKey := fmt.Sprintf("tag:%d", id)
	if cached, found := u.cache.Get(cacheKey); found {
		if t, ok := cached.(*tag.Tag); ok {
			return t, nil
		}
	}

	// Получение тега из репозитория
	t, err := u.tagRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get tag from repository", err)
		return nil, err
	}

	// Сохранение в кеш
	u.cache.Set(cacheKey, t, 5*time.Minute)

	return t, nil
}

// GetTagByName получает тег по имени
func (u *tagUsecase) GetTagByName(name string) (*tag.Tag, error) {
	// Попытка получить тег из кеша
	cacheKey := fmt.Sprintf("tag:name:%s", name)
	if cached, found := u.cache.Get(cacheKey); found {
		if t, ok := cached.(*tag.Tag); ok {
			return t, nil
		}
	}

	// Получение тега из репозитория
	t, err := u.tagRepo.GetByName(name)
	if err != nil {
		logger.Error("Failed to get tag from repository", err)
		return nil, err
	}

	// Сохранение в кеш
	u.cache.Set(cacheKey, t, 5*time.Minute)

	return t, nil
}

// GetPopularTags получает список популярных тегов
func (u *tagUsecase) GetPopularTags(limit int) ([]*tag.Tag, error) {
	// Попытка получить теги из кеша
	cacheKey := fmt.Sprintf("popular_tags:%d", limit)
	if cached, found := u.cache.Get(cacheKey); found {
		if tags, ok := cached.([]*tag.Tag); ok {
			return tags, nil
		}
	}

	// Получение тегов из репозитория
	tags, err := u.tagRepo.GetPopular(limit)
	if err != nil {
		logger.Error("Failed to get popular tags from repository", err)
		return nil, err
	}

	// Сохранение в кеш
	u.cache.Set(cacheKey, tags, 5*time.Minute)

	return tags, nil
}
