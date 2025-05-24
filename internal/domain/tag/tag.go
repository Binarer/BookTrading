package tag

// TagWithCount представляет тег с количеством книг
// Используется для популярных тегов
// Tag - сам тег, BookCount - количество книг с этим тегом
// json:tag для вложенного тега, чтобы структура была удобна для API

type TagWithCount struct {
	Tag       *Tag  `json:"tag"`
	BookCount int64 `json:"book_count"`
}
