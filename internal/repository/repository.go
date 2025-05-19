package repository

import (
	"booktrading/internal/config"
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/state"
	"booktrading/internal/domain/tag"
	"booktrading/internal/domain/user"
	"database/sql"
	"strconv"
)

// BookRepository определяет интерфейс для работы с книгами
type BookRepository interface {
	Create(book *book.Book) error
	GetByID(id uint) (*book.Book, error)
	GetByTags(tagIDs []uint) ([]*book.Book, error)
	AddTags(bookID uint, tagIDs []uint) error
	Update(book *book.Book) error
	Delete(id uint) error
}

// TagRepository определяет интерфейс для работы с тегами
type TagRepository interface {
	Create(tag *tag.Tag) error
	GetByID(id uint) (*tag.Tag, error)
	GetByName(name string) (*tag.Tag, error)
	GetAll() ([]*tag.Tag, error)
	GetPopular(limit int) ([]*tag.Tag, error)
	Update(tag *tag.Tag) error
	Delete(id uint) error
}

// StateRepository определяет интерфейс для работы с состояниями
type StateRepository interface {
	Create(s *state.State) error
	GetByID(id uint) (*state.State, error)
	GetAll() ([]*state.State, error)
	Update(s *state.State) error
	Delete(id uint) error
}

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *user.User) error
	GetByID(id uint) (*user.User, error)
	GetByLogin(login string) (*user.User, error)
	Update(user *user.User) error
	Delete(id uint) error
}

// NewMySQLConnection создает новое подключение к MySQL
func NewMySQLConnection(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + ")/" + cfg.DBName + "?parseTime=true"
	return sql.Open("mysql", dsn)
}
