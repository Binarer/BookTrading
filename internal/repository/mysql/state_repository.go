package mysql

import (
	"booktrading/internal/domain/state"
	"book
	"fmt"
	"gorm.io/gorm"
)

type StateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) *StateRepository {
	return &StateRepository{db: db}
}

func (r *StateRepository) Create(s *state.State) error {
	logger.Info("Creating state in repository with name: " + s.Name)

	if err := r.db.Create(s).Error; err != nil {
		logger.Error("Failed to create state in database", err)
		return err
	}

	logger.Info("State created successfully in database with ID: " + fmt.Sprintf("%d", s.ID))
	return nil
}

func (r *StateRepository) GetByID(id uint) (*state.State, error) {
	var s state.State
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StateRepository) GetAll() ([]*state.State, error) {
	var states []*state.State
	if err := r.db.Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

func (r *StateRepository) Update(s *state.State) error {
	return r.db.Save(s).Error
}

func (r *StateRepository) Delete(id uint) error {
	return r.db.Delete(&state.State{}, id).Error
}
