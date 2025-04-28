package mysql

import (
	"booktrading/internal/domain/tag"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(t *tag.Tag) error {
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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
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

func (r *TagRepository) Update(t *tag.Tag) error {
	return r.db.Save(t).Error
}

func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&tag.Tag{}, id).Error
}

type TagWithCount struct {
	tag.Tag
	Count int
}

func (r *TagRepository) GetPopular(limit int) ([]*tag.Tag, error) {
	var tagCounts []TagWithCount
	if err := r.db.Model(&tag.Tag{}).
		Select("tags.*, COUNT(book_tags.book_id) as count").
		Joins("LEFT JOIN book_tags ON book_tags.tag_id = tags.id").
		Group("tags.id").
		Order("count DESC").
		Limit(limit).
		Find(&tagCounts).Error; err != nil {
		return nil, err
	}

	tags := make([]*tag.Tag, len(tagCounts))
	for i, tc := range tagCounts {
		tags[i] = &tag.Tag{
			Base: tc.Base,
			Name: tc.Name,
		}
	}
	return tags, nil
} 