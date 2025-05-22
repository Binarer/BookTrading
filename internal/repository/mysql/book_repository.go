package mysql

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/tag"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(b *book.Book) error {
	return r.db.Create(b).Error
}

func (r *BookRepository) GetByID(id uint) (*book.Book, error) {
	var b book.Book
	if err := r.db.Preload("Tags").Preload("State").First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookRepository) Update(b *book.Book) error {
	return r.db.Save(b).Error
}

func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&book.Book{}, id).Error
}

func (r *BookRepository) List() ([]*book.Book, error) {
	var books []*book.Book
	if err := r.db.Preload("Tags").Preload("State").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) GetByTag(tagID uint) ([]*book.Book, error) {
	var books []*book.Book
	if err := r.db.Preload("Tags").Preload("State").
		Joins("JOIN book_tags ON book_tags.book_id = books.id").
		Where("book_tags.tag_id = ?", tagID).
		Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) GetByTags(tagIDs []uint) ([]*book.Book, error) {
	var books []*book.Book
	if err := r.db.Preload("Tags").Preload("State").
		Joins("JOIN book_tags ON book_tags.book_id = books.id").
		Where("book_tags.tag_id IN ?", tagIDs).
		Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) AddTags(bookID uint, tagIDs []uint) error {
	// Получаем книгу
	var b book.Book
	if err := r.db.First(&b, bookID).Error; err != nil {
		return err
	}

	// Получаем теги
	var tags []tag.Tag
	if err := r.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}

	// Добавляем теги к книге
	if err := r.db.Model(&b).Association("Tags").Append(tags); err != nil {
		return err
	}

	return nil
}

func (r *BookRepository) GetAll() ([]*book.Book, error) {
	var books []*book.Book
	if err := r.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) GetUserBooks(userID uint, page, pageSize int) ([]*book.Book, int64, error) {
	var books []*book.Book
	var total int64

	// Получаем общее количество книг пользователя
	if err := r.db.Model(&book.Book{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Получаем книги с пагинацией
	offset := (page - 1) * pageSize
	if err := r.db.Preload("Tags").Preload("State").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}
