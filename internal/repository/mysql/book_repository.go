package mysql

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/pkg/logger"
	"database/sql"
	"encoding/json"
)

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *bookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(b *book.Book) error {
	photosJSON, err := json.Marshal(b.Photos)
	if err != nil {
		logger.Error("Failed to marshal photos to JSON", err)
		return err
	}

	query := `
		INSERT INTO books (title, author, description, state_id, photos)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, b.Title, b.Author, b.Description, b.StateID, photosJSON)
	if err != nil {
		logger.Error("Failed to create book in database", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Failed to get last insert ID", err)
		return err
	}
	b.ID = id

	return nil
}

func (r *bookRepository) GetByID(id int64) (*book.Book, error) {
	query := `
		SELECT id, title, author, description, state_id, photos, created_at, updated_at
		FROM books
		WHERE id = ?
	`
	b := &book.Book{}
	var photosJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&b.ID, &b.Title, &b.Author, &b.Description, &b.StateID, &photosJSON,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(photosJSON, &b.Photos); err != nil {
		logger.Error("Failed to unmarshal photos from JSON", err)
		return nil, err
	}

	return b, nil
}

func (r *bookRepository) GetByTags(tagIDs []int64) ([]*book.Book, error) {
	query := `SELECT DISTINCT b.id, b.title, b.author, b.description, b.state_id, b.photos, b.created_at, b.updated_at
			  FROM books b
			  JOIN book_tags bt ON b.id = bt.book_id
			  WHERE bt.tag_id IN (?)`

	rows, err := r.db.Query(query, tagIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*book.Book
	for rows.Next() {
		b := &book.Book{}
		var photosJSON []byte
		err := rows.Scan(
			&b.ID, &b.Title, &b.Author, &b.Description, &b.StateID, &photosJSON,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(photosJSON, &b.Photos); err != nil {
			logger.Error("Failed to unmarshal photos from JSON", err)
			return nil, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (r *bookRepository) AddTags(bookID int64, tagIDs []int64) error {
	query := `
		INSERT INTO book_tags (book_id, tag_id)
		VALUES (?, ?)
	`
	for _, tagID := range tagIDs {
		_, err := r.db.Exec(query, bookID, tagID)
		if err != nil {
			logger.Error("Failed to add tag to book", err)
			return err
		}
	}
	return nil
}

func (r *bookRepository) Update(b *book.Book) error {
	photosJSON, err := json.Marshal(b.Photos)
	if err != nil {
		logger.Error("Failed to marshal photos to JSON", err)
		return err
	}

	query := `
		UPDATE books
		SET title = ?, author = ?, description = ?, state_id = ?, photos = ?
		WHERE id = ?
	`
	_, err = r.db.Exec(query, b.Title, b.Author, b.Description, b.StateID, photosJSON, b.ID)
	return err
}

func (r *bookRepository) Delete(id int64) error {
	query := `
		DELETE FROM books
		WHERE id = ?
	`
	_, err := r.db.Exec(query, id)
	return err
}
