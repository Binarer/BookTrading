package booktag

// BookTag represents a relationship between Book and Tag
type BookTag struct {
	BookID int64 `json:"book_id"`
	TagID  int64 `json:"tag_id"`
} 