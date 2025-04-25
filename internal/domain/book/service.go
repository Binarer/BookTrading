package book

import (
	"booktrading/internal/domain/tag"
	"time"
)

// Service provides business logic for Book entity
type Service struct{}

// NewService creates a new Book service
func NewService() *Service {
	return &Service{}
}

// CreateBook creates a new book with the given details
func (s *Service) CreateBook(title, author, description string) *Book {
	now := time.Now()
	return &Book{
		Title:       title,
		Author:      author,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddTags adds tags to a book
func (s *Service) AddTags(book *Book, tags []*tag.Tag) {
	book.Tags = tags
}

// UpdateBook updates book details
func (s *Service) UpdateBook(book *Book, title, author, description string) {
	book.Title = title
	book.Author = author
	book.Description = description
	book.UpdatedAt = time.Now()
}
