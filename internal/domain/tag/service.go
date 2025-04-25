package tag

import "time"

// Service provides business logic for Tag entity
type Service struct{}

// NewService creates a new Tag service
func NewService() *Service {
	return &Service{}
}

// CreateTag creates a new tag with the given name
func (s *Service) CreateTag(name string) *Tag {
	now := time.Now()
	return &Tag{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateTag updates tag name
func (s *Service) UpdateTag(tag *Tag, name string) {
	tag.Name = name
	tag.UpdatedAt = time.Now()
} 