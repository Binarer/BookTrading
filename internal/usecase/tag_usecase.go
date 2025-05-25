package usecase

import (
	"booktrading/internal/domain/repository"
	"booktrading/internal/domain/tag"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/pkg/logger"
	_ "encoding/json"
	"fmt"
	"time"
)

// TagUseCase определяет интерфейс для работы с тегами
type TagUseCase interface {
	CreateTag(tag *tag.Tag) error
	GetTagByID(id uint) (*tag.Tag, error)
	GetTagByName(name string) (*tag.Tag, error)
	GetAllTags() ([]*tag.Tag, error)
	GetPopularTags(limit int) ([]*tag.TagWithCount, error)
	UpdateTag(id uint, dto *tag.UpdateTagDTO) (*tag.Tag, error)
	DeleteTag(id uint) error
}

// tagUseCase реализует интерфейс TagUseCase
type tagUseCase struct {
	tagRepo  repository.TagRepository
	bookRepo repository.BookRepository
	cache    *cache.Cache
}

// NewTagUseCase создает новый экземпляр TagUseCase
func NewTagUseCase(tagRepo repository.TagRepository, bookRepo repository.BookRepository, cache *cache.Cache) TagUseCase {
	return &tagUseCase{
		tagRepo:  tagRepo,
		bookRepo: bookRepo,
		cache:    cache,
	}
}

// CreateTag создает новый тег
func (u *tagUseCase) CreateTag(t *tag.Tag) error {
	// Check if tag with same name exists
	existingTag, err := u.tagRepo.GetByName(t.Name)
	if err != nil {
		return fmt.Errorf("failed to check tag existence: %w", err)
	}
	if existingTag != nil {
		return fmt.Errorf("tag with name %s already exists", t.Name)
	}

	// Create tag
	if err := u.tagRepo.Create(t); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	// Инвалидируем кеш
	u.cache.DeletePattern("tags:")
	u.cache.Delete("tags:all")

	return nil
}

// GetTagByID получает тег по ID
func (u *tagUseCase) GetTagByID(id uint) (*tag.Tag, error) {
	// Попытка получить тег из кеша
	cacheKey := fmt.Sprintf("tags:id:%d", id)
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
func (u *tagUseCase) GetTagByName(name string) (*tag.Tag, error) {
	// Попытка получить тег из кеша
	cacheKey := fmt.Sprintf("tags:name:%s", name)
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

// GetAllTags получает список всех тегов
func (u *tagUseCase) GetAllTags() ([]*tag.Tag, error) {
	// Попытка получить теги из кеша
	cacheKey := "tags:all"
	if cached, found := u.cache.Get(cacheKey); found {
		if tags, ok := cached.([]*tag.Tag); ok {
			return tags, nil
		}
	}

	// Получение тегов из репозитория
	tags, err := u.tagRepo.GetAll()
	if err != nil {
		logger.Error("Failed to get tags from repository", err)
		return nil, err
	}

	// Сохранение в кеш
	u.cache.Set(cacheKey, tags, 5*time.Minute)

	return tags, nil
}

// GetPopularTags получает список популярных тегов
func (u *tagUseCase) GetPopularTags(limit int) ([]*tag.TagWithCount, error) {
	// Попытка получить теги из кеша
	cacheKey := fmt.Sprintf("tags:popular:%d", limit)
	if cached, found := u.cache.Get(cacheKey); found {
		if tags, ok := cached.([]*tag.TagWithCount); ok {
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

// UpdateTag обновляет существующий тег
func (u *tagUseCase) UpdateTag(id uint, dto *tag.UpdateTagDTO) (*tag.Tag, error) {
	// Get existing tag
	existingTag, err := u.tagRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	if existingTag == nil {
		return nil, fmt.Errorf("tag not found")
	}

	// Update fields
	if dto.Name != "" {
		// Check if new name is already taken by another tag
		if existingTag.Name != dto.Name {
			tagWithName, err := u.tagRepo.GetByName(dto.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to check tag name: %w", err)
			}
			if tagWithName != nil && tagWithName.ID != id {
				return nil, fmt.Errorf("tag with name %s already exists", dto.Name)
			}
		}
		existingTag.Name = dto.Name
	}

	// Update photo if provided
	if dto.Photo != "" {
		existingTag.Photo = dto.Photo
	}

	// Save changes
	if err := u.tagRepo.Update(existingTag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	// Инвалидация кеша
	u.cache.DeletePattern("tags:")
	u.cache.Delete("tags:all")

	return existingTag, nil
}

// DeleteTag удаляет тег по ID
func (u *tagUseCase) DeleteTag(id uint) error {
	// Проверяем, используется ли тег в книгах
	books, err := u.bookRepo.GetByTags([]uint{id})
	if err != nil {
		return err
	}
	if len(books) > 0 {
		return fmt.Errorf("cannot delete tag: it is being used by books")
	}

	// Удаляем тег
	if err := u.tagRepo.Delete(id); err != nil {
		return err
	}

	// Инвалидация кеша
	u.cache.DeletePattern("tags:")
	u.cache.Delete("tags:all")

	return nil
}
