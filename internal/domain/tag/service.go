package tag

// Service предоставляет бизнес-логику для сущности Tag
type Service struct{}

// NewService создает новый сервис для работы с тегами
func NewService() *Service {
	return &Service{}
}

// CreateTag создает новый тег с заданным именем
func (s *Service) CreateTag(name string) *Tag {
	return &Tag{
		Name: name,
	}
}

// UpdateTag обновляет имя тега
func (s *Service) UpdateTag(tag *Tag, name string) {
	tag.Name = name
} 