package booktag

// Service предоставляет бизнес-логику для сущности BookTag
type Service struct{}

// NewService создает новый сервис для работы с связями книг и тегов
func NewService() *Service {
	return &Service{}
}

// CreateBookTag создает новую связь между книгой и тегом
func (s *Service) CreateBookTag(bookID, tagID int64) *BookTag {
	return &BookTag{
		BookID: bookID,
		TagID:  tagID,
	}
}

// ValidateBookTag проверяет, является ли связь валидной
func (s *Service) ValidateBookTag(bookTag *BookTag) bool {
	return bookTag.BookID > 0 && bookTag.TagID > 0
} 