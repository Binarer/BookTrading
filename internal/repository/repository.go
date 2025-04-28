package repository

import (
	"booktrading/internal/config"
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/tag"
	"database/sql"
	"strconv"
)

type BookRepository interface {
	Create(book *book.Book) error
	GetByID(id int64) (*book.Book, error)
	GetByTags(tagIDs []int64) ([]*book.Book, error)
	AddTags(bookID int64, tagIDs []int64) error
	Update(book *book.Book) error
}

type TagRepository interface {
	Create(tag *tag.Tag) error
	GetByID(id int64) (*tag.Tag, error)
	GetByName(name string) (*tag.Tag, error)
	GetPopular(limit int) ([]*tag.Tag, error)
	Update(tag *tag.Tag) error
}

func NewMySQLConnection(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + ")/" + cfg.DBName + "?parseTime=true"
	return sql.Open("mysql", dsn)
}
