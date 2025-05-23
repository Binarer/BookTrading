package usecase

import (
	"booktrading/internal/domain/tag"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/pkg/logger"
	"booktrading/internal/repository"
	"database/sql"
	_ "encoding/json"
	"fmt"
	"time"
)

// TagUsecase определяет интерфейс для работы с тегами
type TagUsecase interface {
	CreateTag(tag *tag.Tag) error
	GetTagByID(id uint) (*tag.Tag, error)
	GetTagByName(name string) (*tag.Tag, error)
	GetAllTags() ([]*tag.Tag, error)
	GetPopularTags(limit int) ([]*tag.TagWithCount, error)
	UpdateTag(id uint, dto *tag.UpdateTagDTO) (*tag.Tag, error)
	DeleteTag(id uint) error
}

// tagUsecase реализует интерфейс TagUsecase
type tagUsecase struct {
	tagRepo  repository.TagRepository
	bookRepo repository.BookRepository
	cache    *cache.Cache
}

// NewTagUsecase создает новый экземпляр tagUsecase
func NewTagUsecase(tagRepo repository.TagRepository, bookRepo repository.BookRepository, cache *cache.Cache) TagUsecase {
	return &tagUsecase{
		tagRepo:  tagRepo,
		bookRepo: bookRepo,
		cache:    cache,
	}
}

// CreateTag создает новый тег
func (u *tagUsecase) CreateTag(t *tag.Tag) error {
	// Создаем тег
	if err := u.tagRepo.Create(t); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	// Инвалидируем кеш
	u.cache.Delete("tags:all")
	u.cache.Delete("tags:popular")

	return nil
}

// GetTagByID получает тег по ID
func (u *tagUsecase) GetTagByID(id uint) (*tag.Tag, error) {
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

// GetAllTags получает список всех тегов
func (u *tagUsecase) GetAllTags() ([]*tag.Tag, error) {
	// Попытка получить теги из кеша
	cacheKey := "all_tags"
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
func (u *tagUsecase) GetPopularTags(limit int) ([]*tag.TagWithCount, error) {
	// Попытка получить теги из кеша
	cacheKey := fmt.Sprintf("popular_tags:%d", limit)
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
func (u *tagUsecase) UpdateTag(id uint, dto *tag.UpdateTagDTO) (*tag.Tag, error) {
	// Получаем существующий тег
	existingTag, err := u.GetTagByID(id)
	if err != nil {
		return nil, err
	}

	// Если имя изменилось, проверяем уникальность
	if dto.Name != "" && dto.Name != existingTag.Name {
		duplicateTag, err := u.tagRepo.GetByName(dto.Name)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if duplicateTag != nil {
			return nil, fmt.Errorf("tag name '%s' is not unique", dto.Name)
		}
		existingTag.Name = dto.Name
	}

	// Обновляем в репозитории
	if err := u.tagRepo.Update(existingTag); err != nil {
		return nil, err
	}

	// Инвалидация кеша
	u.cache.Delete(fmt.Sprintf("tag:%d", id))
	u.cache.Delete(fmt.Sprintf("tag:name:%s", existingTag.Name))

	return existingTag, nil
}

// DeleteTag удаляет тег по ID
func (u *tagUsecase) DeleteTag(id uint) error {
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
	u.cache.Delete("tags")
	u.cache.Delete(fmt.Sprintf("tag:%d", id))

	return nil
}
