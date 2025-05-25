package usecase

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/tag"
	"booktrading/internal/pkg/cache"
	"booktrading/internal/repository/mysql"
	"fmt"
	"time"
)

// BookUseCase определяет интерфейс для работы с книгами
type BookUseCase interface {
	CreateBook(book *book.Book, tagIDs []uint) error
	GetBookByID(id uint) (*book.Book, error)
	GetAllBooks(page, pageSize int) ([]*book.Book, int64, error)
	GetBooksByTags(tagIDs []uint) ([]*book.Book, error)
	AddTagsToBook(bookID uint, tagIDs []uint) error
	UpdateBook(id uint, dto *book.UpdateBookDTO) (*book.Book, error)
	UpdateBookState(id uint, stateID uint) (*book.Book, error)
	DeleteBook(id uint) error
	GetUserBooks(userID uint, page, pageSize int) ([]*book.Book, int64, error)
}

// bookUseCase реализует интерфейс BookUseCase
type bookUseCase struct {
	bookRepo  *mysql.BookRepository
	tagRepo   *mysql.TagRepository
	stateRepo *mysql.StateRepository
	cache     *cache.Cache
	bookSvc   *book.Service
}

// NewBookUseCase создает новый экземпляр bookUseCase
func NewBookUseCase(bookRepo *mysql.BookRepository, tagRepo *mysql.TagRepository, stateRepo *mysql.StateRepository, cache *cache.Cache) BookUseCase {
	return &bookUseCase{
		bookRepo:  bookRepo,
		tagRepo:   tagRepo,
		stateRepo: stateRepo,
		cache:     cache,
		bookSvc:   book.NewService(),
	}
}

// CreateBook создает новую книгу
func (u *bookUseCase) CreateBook(book *book.Book, tagIDs []uint) error {
	// Если состояние не указано, устанавливаем состояние "available"
	if book.StateID == 0 {
		// Используем ID 1 для состояния "available"
		book.StateID = 1
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
	u.cache.DeletePattern("books:")
	u.cache.Delete("books:all")

	return nil
}

// GetBookByID получает книгу по ID
func (u *bookUseCase) GetBookByID(id uint) (*book.Book, error) {
	// Попытка получить книгу из кеша
	cacheKey := fmt.Sprintf("books:id:%d", id)
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
func (u *bookUseCase) GetBooksByTags(tagIDs []uint) ([]*book.Book, error) {
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
func (u *bookUseCase) AddTagsToBook(bookID uint, tagIDs []uint) error {
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
	u.cache.DeletePattern("books:")
	u.cache.Delete("books:all")

	return nil
}

// UpdateBook обновляет существующую книгу
func (u *bookUseCase) UpdateBook(id uint, dto *book.UpdateBookDTO) (*book.Book, error) {
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
		// Проверяем существование нового состояния
		if _, err := u.stateRepo.GetByID(uint(dto.StateID)); err != nil {
			return nil, fmt.Errorf("invalid state ID: %w", err)
		}
		existingBook.StateID = uint(dto.StateID)
	}

	// Обновляем в репозитории
	if err := u.bookRepo.Update(existingBook); err != nil {
		return nil, err
	}

	// Инвалидация кеша
	u.cache.DeletePattern("books:")
	u.cache.Delete("books:all")

	return existingBook, nil
}

// UpdateBookState обновляет состояние книги
func (u *bookUseCase) UpdateBookState(id uint, stateID uint) (*book.Book, error) {
	// Получаем существующую книгу
	existingBook, err := u.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	// Проверяем существование нового состояния
	if _, err := u.stateRepo.GetByID(stateID); err != nil {
		return nil, fmt.Errorf("invalid state ID: %w", err)
	}

	// Обновляем состояние
	existingBook.StateID = stateID

	// Обновляем в репозитории
	if err := u.bookRepo.Update(existingBook); err != nil {
		return nil, err
	}

	// Инвалидация кеша
	u.cache.DeletePattern("books:")
	u.cache.Delete("books:all")

	return existingBook, nil
}

// DeleteBook удаляет книгу
func (u *bookUseCase) DeleteBook(id uint) error {
	// Удаляем из репозитория
	if err := u.bookRepo.Delete(id); err != nil {
		return err
	}

	// Инвалидация кеша
	u.cache.DeletePattern("books:")
	u.cache.Delete("books:all")

	return nil
}

// GetAllBooks получает все книги с пагинацией
func (u *bookUseCase) GetAllBooks(page, pageSize int) ([]*book.Book, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Попытка получить книги из кеша
	cacheKey := fmt.Sprintf("books:page:%d:size:%d", page, pageSize)
	if cached, found := u.cache.Get(cacheKey); found {
		if result, ok := cached.(map[string]interface{}); ok {
			if books, ok := result["books"].([]*book.Book); ok {
				if total, ok := result["total"].(int64); ok {
					return books, total, nil
				}
			}
		}
	}

	// Получение книг из репозитория
	books, total, err := u.bookRepo.GetAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Сохранение в кеш
	result := map[string]interface{}{
		"books": books,
		"total": total,
	}
	u.cache.Set(cacheKey, result, 5*time.Minute)

	return books, total, nil
}

// GetUserBooks получает книги пользователя с пагинацией
func (u *bookUseCase) GetUserBooks(userID uint, page, pageSize int) ([]*book.Book, int64, error) {
	// Валидация параметров пагинации
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return u.bookRepo.GetUserBooks(userID, page, pageSize)
}
