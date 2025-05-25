package mysql

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/repository"
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/domain/user"
	"booktrading/internal/pkg/logger"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) repository.BookRepository {
	return &BookRepository{db: db}
}

// validatePhotos проверяет фотографии на соответствие требованиям
func validatePhotos(photos []string) error {
	if len(photos) > 5 {
		return errors.New("maximum 5 photos allowed")
	}

	for _, photo := range photos {
		// Проверяем формат base64
		if !strings.HasPrefix(photo, "data:image/") {
			return errors.New("invalid photo format: must be base64 encoded image")
		}

		// Извлекаем base64 данные
		parts := strings.Split(photo, ",")
		if len(parts) != 2 {
			return errors.New("invalid photo format: missing base64 data")
		}

		// Проверяем размер (5MB = 5 * 1024 * 1024 байт)
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return fmt.Errorf("invalid base64 data: %v", err)
		}

		if len(decoded) > 5*1024*1024 {
			return errors.New("photo size exceeds 5MB limit")
		}

		// Проверяем формат изображения
		contentType := strings.TrimPrefix(parts[0], "data:")
		if !strings.HasPrefix(contentType, "image/jpeg") && !strings.HasPrefix(contentType, "image/png") {
			return errors.New("unsupported image format: only JPEG and PNG are allowed")
		}
	}

	return nil
}

func (r *BookRepository) Create(bookData *book.Book) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Проверяем существование пользователя
	var userData user.User
	if err := tx.First(&userData, bookData.UserID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("user not found: %w", err)
	}

	// Проверяем уникальность названия книги
	var existingBookData book.Book
	if err := tx.Where("title = ?", bookData.Title).First(&existingBookData).Error; err == nil {
		tx.Rollback()
		return fmt.Errorf("book with title '%s' already exists", bookData.Title)
	}

	// Проверяем существование всех тегов
	for _, tagData := range bookData.Tags {
		var existingTagData tag.Tag
		if err := tx.First(&existingTagData, tagData.ID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("tag with ID %d not found", tagData.ID)
		}
	}

	// Проверяем существование состояния
	var existingStateData state.State
	if err := tx.First(&existingStateData, bookData.StateID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("state with ID %d not found", bookData.StateID)
	}

	// Проверяем фотографии только если они есть
	if len(bookData.Photos) > 0 {
		photoURLs := make([]string, len(bookData.Photos))
		for i, photo := range bookData.Photos {
			// Проверяем, что photo_url содержит base64 строку
			if !strings.HasPrefix(photo.PhotoURL, "data:image/") {
				tx.Rollback()
				return fmt.Errorf("invalid photo format at index %d: must be base64 encoded image", i)
			}
			photoURLs[i] = photo.PhotoURL
		}
		if err := validatePhotos(photoURLs); err != nil {
			tx.Rollback()
			return fmt.Errorf("invalid photos: %w", err)
		}
	}

	// Сохраняем книгу
	if err := tx.Create(bookData).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create book: %w", err)
	}

	// Создаем связи с тегами
	if err := tx.Model(bookData).Association("Tags").Replace(bookData.Tags); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create book-tag relations: %w", err)
	}

	// Создаем фотографии, если они есть
	if len(bookData.Photos) > 0 {
		for i, photo := range bookData.Photos {
			photo.BookID = bookData.ID
			photo.IsMain = i == 0 // Первая фотография - главная
			// Сбрасываем ID, чтобы GORM создал новый
			photo.ID = 0
			if err := tx.Create(&photo).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create photo: %w", err)
			}
		}
	}

	return tx.Commit().Error
}

func (r *BookRepository) GetByID(id uint) (*book.Book, error) {
	var b book.Book
	if err := r.db.Preload("Photos").First(&b, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *BookRepository) Update(b *book.Book) error {
	// Проверяем существование книги
	var existingBook book.Book
	if err := r.db.Preload("Photos").First(&existingBook, b.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("Book not found", fmt.Errorf("book with ID %d not found", b.ID))
			return errors.New("book not found")
		}
		logger.Error("Failed to get book", err)
		return err
	}

	// Проверяем уникальность названия книги для пользователя
	var count int64
	if err := r.db.Model(&book.Book{}).Where("user_id = ? AND title = ? AND id != ?", b.UserID, b.Title, b.ID).Count(&count).Error; err != nil {
		logger.Error("Failed to check book title uniqueness", err)
		return err
	}
	if count > 0 {
		logger.Error("Book title already exists", fmt.Errorf("book with title '%s' already exists for user %d", b.Title, b.UserID))
		return errors.New("book title already exists")
	}

	// Проверяем существование тегов
	if len(b.Tags) > 0 {
		var tagIDs []uint
		for _, t := range b.Tags {
			tagIDs = append(tagIDs, t.ID)
		}
		var count int64
		if err := r.db.Model(&tag.Tag{}).Where("id IN ?", tagIDs).Count(&count).Error; err != nil {
			logger.Error("Failed to check tags existence", err)
			return err
		}
		if int(count) != len(tagIDs) {
			logger.Error("Some tags not found", fmt.Errorf("some tags from %v not found", tagIDs))
			return errors.New("some tags not found")
		}
	}

	// Проверяем существование состояния
	if b.StateID != 0 {
		var stateExists bool
		if err := r.db.Model(&state.State{}).Select("1").Where("id = ?", b.StateID).Take(&stateExists).Error; err != nil {
			logger.Error("Failed to check state existence", err)
			return err
		}
		if !stateExists {
			logger.Error("State not found", fmt.Errorf("state with ID %d not found", b.StateID))
			return errors.New("state not found")
		}
	}

	// Валидируем фотографии
	var photoURLs []string
	for _, photo := range b.Photos {
		photoURLs = append(photoURLs, photo.PhotoURL)
	}
	if err := validatePhotos(photoURLs); err != nil {
		logger.Error("Invalid photos", err)
		return err
	}

	// Обновляем книгу в транзакции
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Обновляем основные данные книги
		if err := tx.Save(b).Error; err != nil {
			return err
		}

		// Удаляем старые фотографии
		if err := tx.Where("book_id = ?", b.ID).Delete(&book.BookPhoto{}).Error; err != nil {
			return err
		}

		// Создаем новые фотографии
		for i := range b.Photos {
			b.Photos[i].BookID = b.ID
		}
		if err := tx.Create(&b.Photos).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *BookRepository) Delete(id uint) error {
	// Проверяем существование книги
	var b book.Book
	if err := r.db.First(&b, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("Book not found", fmt.Errorf("book with ID %d not found", id))
			return errors.New("book not found")
		}
		logger.Error("Failed to check book existence", err)
		return err
	}

	// Удаляем книгу в транзакции
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Удаляем фотографии
		if err := tx.Where("book_id = ?", id).Delete(&book.BookPhoto{}).Error; err != nil {
			return err
		}

		// Удаляем книгу
		if err := tx.Delete(&book.Book{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *BookRepository) List() ([]*book.Book, error) {
	var books []*book.Book
	if err := r.db.Preload("Tags").Preload("State").Preload("User").Preload("Photos").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) GetByTag(tagID uint) ([]*book.Book, error) {
	var books []*book.Book
	if err := r.db.Preload("Tags").Preload("State").Preload("User").
		Joins("JOIN book_tags ON book_tags.book_id = books.id").
		Where("book_tags.tag_id = ?", tagID).
		Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) GetByTags(tagIDs []uint) ([]*book.Book, error) {
	var books []*book.Book
	if err := r.db.Preload("Tags").Preload("State").Preload("User").Preload("Photos").
		Joins("JOIN book_tags ON books.id = book_tags.book_id").
		Where("book_tags.tag_id IN ?", tagIDs).
		Group("books.id").
		Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) AddTags(bookID uint, tagIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, tagID := range tagIDs {
			if err := tx.Exec("INSERT INTO book_tags (book_id, tag_id) VALUES (?, ?)", bookID, tagID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *BookRepository) GetAll(page, pageSize int) ([]*book.Book, int64, error) {
	var books []*book.Book
	var total int64

	if err := r.db.Model(&book.Book{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Preload("Tags").Preload("State").Preload("User").Preload("Photos").
		Offset(offset).Limit(pageSize).
		Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *BookRepository) GetUserBooks(userID uint, page, pageSize int) ([]*book.Book, int64, error) {
	var books []*book.Book
	var total int64

	if err := r.db.Model(&book.Book{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Preload("Tags").Preload("State").Preload("User").Preload("Photos").
		Where("user_id = ?", userID).
		Offset(offset).Limit(pageSize).
		Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}
