package usecase

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/booktag"
	"booktrading/internal/domain/tag"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/repository"
	"fmt"
	"time"
)

// BookUsecase определяет интерфейс для работы с книгами
type BookUsecase interface {
	CreateBook(book *book.Book, tagIDs []int64) error
	GetBookByID(id int64) (*book.Book, error)
	GetBooksByTags(tagIDs []int64) ([]*book.Book, error)
	AddTagsToBook(bookID int64, tagIDs []int64) error
	UpdateBook(id int64, dto *book.UpdateBookDTO) (*book.Book, error)
	DeleteBook(id int64) error
}

// bookUsecase реализует интерфейс BookUsecase
type bookUsecase struct {
	bookRepo    repository.BookRepository
	tagRepo     repository.TagRepository
	cache       *cache.Cache
	bookSvc     *book.Service
	bookTagSvc  *booktag.Service
}

// NewBookUsecase создает новый экземпляр bookUsecase
func NewBookUsecase(bookRepo repository.BookRepository, tagRepo repository.TagRepository, cache *cache.Cache) BookUsecase {
	return &bookUsecase{
		bookRepo:    bookRepo,
		tagRepo:     tagRepo,
		cache:       cache,
		bookSvc:     book.NewService(),
		bookTagSvc:  booktag.NewService(),
	}
}

// CreateBook создает новую книгу
func (u *bookUsecase) CreateBook(book *book.Book, tagIDs []int64) error {
	// Устанавливаем состояние по умолчанию
	if book.StateID == 0 {
		book.StateID = 1 // Assuming 1 is the default state ID
	}

	// Сохраняем в репозиторий
	if err := u.bookRepo.Create(book); err != nil {
		return err
	}

	// Добавляем теги, если они указаны
	if len(tagIDs) > 0 {
		if err := u.bookRepo.AddTags(book.ID, tagIDs); err != nil {
			return err
		}
	}

	// Инвалидация кеша
	u.cache.Delete("books")

	return nil
}

// GetBookByID получает книгу по ID
func (u *bookUsecase) GetBookByID(id int64) (*book.Book, error) {
	// Попытка получить книгу из кеша
	cacheKey := fmt.Sprintf("book:%d", id)
	if cached, found := u.cache.Get(cacheKey); found {
		if book, ok := cached.(*book.Book); ok {
			return book, nil
		}
	}

	// Получение книги из репозитория
	book, err := u.bookRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Сохранение в кеш
	u.cache.Set(cacheKey, book, 5*time.Minute)

	return book, nil
}

// GetBooksByTags получает книги по тегам
func (u *bookUsecase) GetBooksByTags(tagIDs []int64) ([]*book.Book, error) {
	// Попытка получить книги из кеша
	cacheKey := fmt.Sprintf("books:tags:%v", tagIDs)
	if cached, found := u.cache.Get(cacheKey); found {
		if books, ok := cached.([]*book.Book); ok {
			return books, nil
		}
	}

	// Получение книг из репозитория
	books, err := u.bookRepo.GetByTags(tagIDs)
	if err != nil {
		return nil, err
	}

	// Сохранение в кеш
	u.cache.Set(cacheKey, books, 5*time.Minute)

	return books, nil
}

// AddTagsToBook добавляет теги к книге
func (u *bookUsecase) AddTagsToBook(bookID int64, tagIDs []int64) error {
	// Получаем книгу
	book, err := u.GetBookByID(bookID)
	if err != nil {
		return err
	}

	// Получаем теги
	var tags []*tag.Tag
	for _, tagID := range tagIDs {
		tag, err := u.tagRepo.GetByID(tagID)
		if err != nil {
			return err
		}
		tags = append(tags, tag)
	}

	// Добавляем теги через доменный сервис
	u.bookSvc.AddTags(book, tags)

	// Сохраняем связи в репозитории
	if err := u.bookRepo.AddTags(bookID, tagIDs); err != nil {
		return err
	}

	// Инвалидация кеша
	u.cache.Delete(fmt.Sprintf("book:%d", bookID))
	u.cache.Delete("books")

	return nil
}

// UpdateBook обновляет существующую книгу
func (u *bookUsecase) UpdateBook(id int64, dto *book.UpdateBookDTO) (*book.Book, error) {
	// Получаем существующую книгу
	existingBook, err := u.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	// Обновляем поля
	if dto.Title != "" {
		existingBook.Title = dto.Title
	}
	if dto.Author != "" {
		existingBook.Author = dto.Author
	}
	if dto.Description != "" {
		existingBook.Description = dto.Description
	}
	if dto.StateID != 0 {
		existingBook.StateID = dto.StateID
	}
	if len(dto.Photos) > 0 {
		existingBook.Photos = dto.Photos
	}

	// Обновляем в репозитории
	if err := u.bookRepo.Update(existingBook); err != nil {
		return nil, err
	}

	// Инвалидация кеша
	u.cache.Delete(fmt.Sprintf("book:%d", id))
	u.cache.Delete("books")

	return existingBook, nil
}

// DeleteBook удаляет книгу
func (u *bookUsecase) DeleteBook(id int64) error {
	return u.bookRepo.Delete(id)
} 