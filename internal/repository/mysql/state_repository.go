package mysql

import (
	"booktrading/internal/domain/state"
	"booktrading/internal/pkg/logger"
	"errors"
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

	var count int64
	if err := r.db.Model(&state.State{}).Where("name = ?", s.Name).Count(&count).Error; err != nil {
		logger.Error("Failed to check state name uniqueness", err)
		return err
	}
	if count > 0 {
		return errors.New("state name must be unique")
	}

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
		if err == gorm.ErrRecordNotFound {
			logger.Error("State not found", fmt.Errorf("state with ID %d not found", id))
			return nil, errors.New("state not found")
		}
		logger.Error("Failed to get state by ID", err)
		return nil, err
	}
	return &s, nil
}

func (r *StateRepository) GetAll() ([]*state.State, error) {
	var states []*state.State
	if err := r.db.Find(&states).Error; err != nil {
		logger.Error("Failed to get all states", err)
		return nil, err
	}
	return states, nil
}

func (r *StateRepository) Update(s *state.State) error {
	var count int64
	if err := r.db.Model(&state.State{}).
		Where("name = ? AND id != ?", s.Name, s.ID).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check state name uniqueness", err)
		return err
	}
	if count > 0 {
		return errors.New("state name must be unique")
	}

	if err := r.db.Save(s).Error; err != nil {
		logger.Error("Failed to update state", err)
		return err
	}
	return nil
}

func (r *StateRepository) Delete(id uint) error {
	var count int64
	if err := r.db.Model(&state.State{}).
		Joins("JOIN books ON books.state_id = states.id").
		Where("states.id = ?", id).
		Count(&count).Error; err != nil {
		logger.Error("Failed to check state usage in books", err)
		return err
	}

	if count > 0 {
		return errors.New("cannot delete state: it is used in books")
	}

	if err := r.db.Delete(&state.State{}, id).Error; err != nil {
		logger.Error("Failed to delete state", err)
		return err
	}
	return nil
}
