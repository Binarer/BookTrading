package usecase

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/tag"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/repository/mysql"
	"fmt"
	"time"
)

// BookUsecase определяет интерфейс для работы с книгами
type BookUsecase interface {
	CreateBook(book *book.Book, tagIDs []uint) error
	GetBookByID(id uint) (*book.Book, error)
	GetBooksByTags(tagIDs []uint) ([]*book.Book, error)
	AddTagsToBook(bookID uint, tagIDs []uint) error
	UpdateBook(id uint, dto *book.UpdateBookDTO) (*book.Book, error)
	UpdateBookState(id uint, stateID uint) (*book.Book, error)
	DeleteBook(id uint) error
}

// bookUsecase реализует интерфейс BookUsecase
type bookUsecase struct {
	bookRepo *mysql.BookRepository
	tagRepo  *mysql.TagRepository
	cache    *cache.Cache
	bookSvc  *book.Service
}

// NewBookUsecase создает новый экземпляр bookUsecase
func NewBookUsecase(bookRepo *mysql.BookRepository, tagRepo *mysql.TagRepository, cache *cache.Cache) BookUsecase {
	return &bookUsecase{
		bookRepo: bookRepo,
		tagRepo:  tagRepo,
		cache:    cache,
		bookSvc:  book.NewService(),
	}
}

// CreateBook создает новую книгу
func (u *bookUsecase) CreateBook(book *book.Book, tagIDs []uint) error {
	// Устанавливаем состояние по умолчанию
	if book.StateID == 0 {
		book.StateID = 1 // Предполагаем, что 1 - это ID состояния по умолчанию
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

	// Добавляем теги к книге
	u.bookSvc.AddTags(book, tags)

	// Сохраняем в репозиторий
	if err := u.bookRepo.Create(book); err != nil {
		return err
	}

	// Инвалидация кеша
	u.cache.Delete("books")

	return nil
}

// GetBookByID получает книгу по ID
func (u *bookUsecase) GetBookByID(id uint) (*book.Book, error) {
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
func (u *bookUsecase) GetBooksByTags(tagIDs []uint) ([]*book.Book, error) {
	// Попытка получить книги из кеша
	cacheKey := fmt.Sprintf("books:tags:%v", tagIDs)
	if cached, found := u.cache.Get(cacheKey); found {
		if books, ok := cached.([]*book.Book); ok {
			return books, nil
		}
	}

	// Получение книг из репозитория
	var books []*book.Book
	for _, tagID := range tagIDs {
		tagBooks, err := u.bookRepo.GetByTag(tagID)
		if err != nil {
			return nil, err
		}
		books = append(books, tagBooks...)
	}

	// Сохранение в кеш
	u.cache.Set(cacheKey, books, 5*time.Minute)

	return books, nil
}

// AddTagsToBook добавляет теги к книге
func (u *bookUsecase) AddTagsToBook(bookID uint, tagIDs []uint) error {
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

	// Сохраняем в репозитории
	if err := u.bookRepo.Update(book); err != nil {
		return err
	}

	// Инвалидация кеша
	u.cache.Delete(fmt.Sprintf("book:%d", bookID))
	u.cache.Delete("books")

	return nil
}

// UpdateBook обновляет существующую книгу
func (u *bookUsecase) UpdateBook(id uint, dto *book.UpdateBookDTO) (*book.Book, error) {
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
		existingBook.StateID = uint(dto.StateID)
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

// UpdateBookState обновляет состояние книги
func (u *bookUsecase) UpdateBookState(id uint, stateID uint) (*book.Book, error) {
	// Получаем существующую книгу
	existingBook, err := u.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	// Обновляем состояние
	existingBook.StateID = stateID

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
func (u *bookUsecase) DeleteBook(id uint) error {
	// Удаляем из репозитория
	if err := u.bookRepo.Delete(id); err != nil {
		return err
	}

	// Инвалидация кеша
	u.cache.Delete(fmt.Sprintf("book:%d", id))
	u.cache.Delete("books")

	return nil
} 