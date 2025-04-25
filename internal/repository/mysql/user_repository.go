package mysql

import (
	"booktrading/internal/domain/user"
	"database/sql"
	"time"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(u *user.User) error {
	query := `INSERT INTO users (username, email, password_hash, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, u.Username, u.Email, u.PasswordHash, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = id
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

func (r *userRepository) GetByID(id int64) (*user.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at 
			  FROM users WHERE id = ?`

	u := &user.User{}
	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userRepository) GetByEmail(email string) (*user.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at 
			  FROM users WHERE email = ?`

	u := &user.User{}
	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userRepository) Update(u *user.User) error {
	query := `UPDATE users 
			  SET username = ?, email = ?, password_hash = ?, updated_at = ? 
			  WHERE id = ?`

	now := time.Now()
	_, err := r.db.Exec(query, u.Username, u.Email, u.PasswordHash, now, u.ID)
	if err != nil {
		return err
	}

	u.UpdatedAt = now
	return nil
} 