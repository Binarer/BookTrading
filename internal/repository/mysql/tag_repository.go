package mysql

import (
	"booktrading/internal/domain/tag"
	"database/sql"
	"time"
)

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *tagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(t *tag.Tag) error {
	query := `INSERT INTO tags (name, created_at, updated_at) 
			  VALUES (?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, t.Name, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = id
	t.CreatedAt = now
	t.UpdatedAt = now
	return nil
}

func (r *tagRepository) GetByID(id int64) (*tag.Tag, error) {
	query := `SELECT id, name, created_at, updated_at 
			  FROM tags WHERE id = ?`

	t := &tag.Tag{}
	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *tagRepository) GetByName(name string) (*tag.Tag, error) {
	query := `SELECT id, name, created_at, updated_at 
			  FROM tags WHERE name = ?`

	t := &tag.Tag{}
	err := r.db.QueryRow(query, name).Scan(
		&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *tagRepository) GetPopular(limit int) ([]*tag.Tag, error) {
	query := `SELECT t.id, t.name, t.created_at, t.updated_at, COUNT(bt.book_id) as usage_count
			  FROM tags t
			  LEFT JOIN book_tags bt ON t.id = bt.tag_id
			  GROUP BY t.id
			  ORDER BY usage_count DESC
			  LIMIT ?`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*tag.Tag
	for rows.Next() {
		t := &tag.Tag{}
		var usageCount int
		err := rows.Scan(
			&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt, &usageCount,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}

func (r *tagRepository) Update(t *tag.Tag) error {
	query := `UPDATE tags 
			  SET name = ?, updated_at = ?
			  WHERE id = ?`

	now := time.Now()
	_, err := r.db.Exec(query, t.Name, now, t.ID)
	if err != nil {
		return err
	}

	t.UpdatedAt = now
	return nil
}

// Delete удаляет тег по ID
func (r *tagRepository) Delete(id int64) error {
	query := `DELETE FROM tags WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
} 