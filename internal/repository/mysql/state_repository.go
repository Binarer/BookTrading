package mysql

import (
	"booktrading/internal/domain/state"
	"booktrading/internal/pkg/logger"
	"database/sql"
)

type stateRepository struct {
	db *sql.DB
}

func NewStateRepository(db *sql.DB) *stateRepository {
	return &stateRepository{db: db}
}

func (r *stateRepository) Create(s *state.State) error {
	query := `
		INSERT INTO states (name)
		VALUES (?)
	`
	result, err := r.db.Exec(query, s.Name)
	if err != nil {
		logger.Error("Failed to create state in database", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Failed to get last insert ID", err)
		return err
	}
	s.ID = id

	return nil
}

func (r *stateRepository) GetByID(id int64) (*state.State, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM states
		WHERE id = ?
	`
	s := &state.State{}
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.Name, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *stateRepository) GetAll() ([]*state.State, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM states
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*state.State
	for rows.Next() {
		s := &state.State{}
		err := rows.Scan(
			&s.ID, &s.Name, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		states = append(states, s)
	}

	return states, nil
}

func (r *stateRepository) Update(s *state.State) error {
	query := `
		UPDATE states
		SET name = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, s.Name, s.ID)
	return err
}

func (r *stateRepository) Delete(id int64) error {
	query := `
		DELETE FROM states
		WHERE id = ?
	`
	_, err := r.db.Exec(query, id)
	return err
} 