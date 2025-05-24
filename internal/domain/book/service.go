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
		StateID:     1, // Состояние "available" по умолчанию
	}
}

// AddTags добавляет теги к книге
func (s *Service) AddTags(book *Book, tags []*tag.Tag) {
	// Используем указатели на теги для сохранения ссылок
	book.Tags = tags
}

// UpdateBook обновляет детали книги
func (s *Service) UpdateBook(book *Book, title, author, description string) {
	if title != "" {
		book.Title = title
	}
	if author != "" {
		book.Author = author
	}
	if description != "" {
		book.Description = description
	}
}
