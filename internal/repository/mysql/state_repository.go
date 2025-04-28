package mysql

import (
	"booktrading/internal/domain/state"
	"gorm.io/gorm"
)

type StateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) *StateRepository {
	return &StateRepository{db: db}
}

func (r *StateRepository) Create(s *state.State) error {
	return r.db.Create(s).Error
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