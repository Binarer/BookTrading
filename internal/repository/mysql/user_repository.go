package mysql

import (
	"booktrading/internal/domain/book"
	"booktrading/internal/domain/user"
	"booktrading/internal/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) GetByID(id uint) (*user.User, error) {
	var u user.User
	if err := r.db.First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	// Загружаем книги пользователя
	var books []*book.Book
	if err := r.db.Where("user_id = ?", id).Find(&books).Error; err != nil {
		return nil, err
	}
	u.Books = books

	return &u, nil
}

func (r *UserRepository) GetByLogin(login string) (*user.User, error) {
	var u user.User
	if err := r.db.Where("login = ?", login).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(u *user.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&user.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *UserRepository) GetAll() ([]*user.User, error) {
	var users []*user.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	// Загружаем книги для каждого пользователя
	for _, u := range users {
		var books []*book.Book
		if err := r.db.Where("user_id = ?", u.ID).Find(&books).Error; err != nil {
			return nil, err
		}
		u.Books = books
	}

	return users, nil
}
