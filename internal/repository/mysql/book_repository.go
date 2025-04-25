package mysql

import (
	"booktrading/internal/domain/book"
	"database/sql"
	"time"
)

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *bookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(b *book.Book) error {
	query := `INSERT INTO books (title, author, description, user_id, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, b.Title, b.Author, b.Description, b.UserID, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	b.ID = id
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (r *bookRepository) GetByID(id int64) (*book.Book, error) {
	query := `SELECT id, title, author, description, user_id, created_at, updated_at 
			  FROM books WHERE id = ?`

	b := &book.Book{}
	err := r.db.QueryRow(query, id).Scan(
		&b.ID, &b.Title, &b.Author, &b.Description, &b.UserID,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *bookRepository) GetByTags(tagIDs []int64) ([]*book.Book, error) {
	query := `SELECT DISTINCT b.id, b.title, b.author, b.description, b.user_id, b.created_at, b.updated_at
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
		err := rows.Scan(
			&b.ID, &b.Title, &b.Author, &b.Description, &b.UserID,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	return books, nil
}

func (r *bookRepository) AddTags(bookID int64, tagIDs []int64) error {
	query := `INSERT INTO book_tags (book_id, tag_id) VALUES (?, ?)`

	for _, tagID := range tagIDs {
		_, err := r.db.Exec(query, bookID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}
