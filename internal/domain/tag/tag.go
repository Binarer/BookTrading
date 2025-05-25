package tag

// TagWithCount представляет тег с количеством книг
// @Description Модель тега с информацией о количестве книг, использующих этот тег
// Tag - сам тег, BookCount - количество книг с этим тегом
// json:tag для вложенного тега, чтобы структура была удобна для API

type TagWithCount struct {
	// @Description Информация о теге
	// @example {"id": 1, "name": "fiction", "photo": "data:image/jpeg;base64,/9j/4AAQSkZJRg..."}
	Tag *Tag `json:"tag"`
	// @Description Количество книг с этим тегом
	// @example 42
	BookCount int64 `json:"book_count"`
}
