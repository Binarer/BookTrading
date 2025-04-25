package booktag

// Service provides business logic for BookTag entity
type Service struct{}

// NewService creates a new BookTag service
func NewService() *Service {
	return &Service{}
}

// CreateBookTag creates a new relationship between book and tag
func (s *Service) CreateBookTag(bookID, tagID int64) *BookTag {
	return &BookTag{
		BookID: bookID,
		TagID:  tagID,
	}
}

// ValidateBookTag checks if the relationship is valid
func (s *Service) ValidateBookTag(bookTag *BookTag) bool {
	return bookTag.BookID > 0 && bookTag.TagID > 0
} 