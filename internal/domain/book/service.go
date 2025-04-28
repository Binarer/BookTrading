package book

import (
	"booktrading/internal/domain/tag"
)

// Service предоставляет бизнес-логику для сущности Book
type Service struct{}

// NewService создает новый сервис для работы с книгами
func NewService() *Service {
	return &Service{}
}

// CreateBook создает новую книгу с заданными деталями
func (s *Service) CreateBook(title, author, description string) *Book {
	return &Book{
		Title:       title,
		Author:      author,
		Description: description,
	}
}

// AddTags добавляет теги к книге
func (s *Service) AddTags(book *Book, tags []*tag.Tag) {
	bookTags := make([]tag.Tag, len(tags))
	for i, t := range tags {
		bookTags[i] = *t
	}
	book.Tags = bookTags
}

// UpdateBook обновляет детали книги
func (s *Service) UpdateBook(book *Book, title, author, description string) {
	book.Title = title
	book.Author = author
	book.Description = description
}
