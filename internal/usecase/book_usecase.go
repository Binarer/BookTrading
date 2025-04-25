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
	CreateBook(book *book.Book) error
	GetBookByID(id int64) (*book.Book, error)
	GetBooksByTags(tagIDs []int64) ([]*book.Book, error)
	AddTagsToBook(bookID int64, tagIDs []int64) error
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
func (u *bookUsecase) CreateBook(book *book.Book) error {
	// Создаем книгу через доменный сервис
	newBook := u.bookSvc.CreateBook(book.Title, book.Author, book.Description)

	// Сохраняем в репозиторий
	if err := u.bookRepo.Create(newBook); err != nil {
		return err
	}

	// Копируем ID в исходный объект
	book.ID = newBook.ID
	book.CreatedAt = newBook.CreatedAt
	book.UpdatedAt = newBook.UpdatedAt

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

	return nil
} 