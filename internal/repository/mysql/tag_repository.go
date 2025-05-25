package mysql

import (
	"booktrading/internal/domain/repository"
	"booktrading/internal/domain/tag"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(t *tag.Tag) error {
	// Проверяем существование тега с таким именем
	existingTag, err := r.GetByName(t.Name)
	if err != nil {
		return fmt.Errorf("failed to check tag existence: %w", err)
	}

	// Если тег уже существует, возвращаем ошибку
	if existingTag != nil {
		return fmt.Errorf("tag with name %s already exists", t.Name)
	}

	// Создаем новый тег
	return r.db.Create(t).Error
}

func (r *TagRepository) GetByID(id uint) (*tag.Tag, error) {
	var t tag.Tag
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TagRepository) GetByName(name string) (*tag.Tag, error) {
	var t tag.Tag
	if err := r.db.Where("name = ?", name).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Возвращаем nil вместо ошибки, если тег не найден
		}
		return nil, err
	}
	return &t, nil
}

func (r *TagRepository) GetAll() ([]*tag.Tag, error) {
	var tags []*tag.Tag
	if err := r.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) GetPopular(limit int) ([]*tag.TagWithCount, error) {
	var results []struct {
		tag.Tag
		BookCount int64
	}
	if err := r.db.Model(&tag.Tag{}).
		Select("tags.*, COUNT(book_tags.book_id) as book_count").
		Joins("LEFT JOIN book_tags ON book_tags.tag_id = tags.id").
		Group("tags.id").
		Order("book_count DESC").
		Limit(limit).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	var tagsWithCount []*tag.TagWithCount
	for _, res := range results {
		t := res.Tag // копия
		tagsWithCount = append(tagsWithCount, &tag.TagWithCount{
			Tag:       &t,
			BookCount: res.BookCount,
		})
	}
	return tagsWithCount, nil
}

func (r *TagRepository) Update(t *tag.Tag) error {
	return r.db.Save(t).Error
}

func (r *TagRepository) Delete(id uint) error {
	// Проверяем, используется ли тег в книгах
	var count int64
	if err := r.db.Model(&tag.Tag{}).
		Joins("JOIN book_tags ON book_tags.tag_id = tags.id").
		Where("tags.id = ?", id).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete tag: it is used in books")
	}

	return r.db.Delete(&tag.Tag{}, id).Error
}
