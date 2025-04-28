package booktag

// BookTag represents a relationship between Book and Tag
// @Description BookTag model for linking books and tags
type BookTag struct {
	// @Description ID of the book
	// @example 1
	BookID int64 `json:"book_id"`

	// @Description ID of the tag
	// @example 1
	TagID int64 `json:"tag_id"`

	// @Description When the relationship was created
	// @example 2025-04-28T12:00:00Z
	CreatedAt string `json:"created_at"`
} 