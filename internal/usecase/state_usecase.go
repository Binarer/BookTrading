package usecase

import (
	"booktrading/internal/domain/repository"
	"booktrading/internal/domain/state"
	"booktrading/internal/pkg/logger"
	"fmt"
)

// StateUsecase определяет интерфейс для работы с состояниями книг
type StateUsecase interface {
	Create(s *state.State) error
	GetByID(id uint) (*state.State, error)
	GetAll() ([]*state.State, error)
	Update(s *state.State) error
	Delete(id uint) error
}

// stateUsecase реализует интерфейс StateUsecase
type stateUsecase struct {
	stateRepo repository.StateRepository
}

// NewStateUsecase создает новый экземпляр stateUsecase
func NewStateUsecase(stateRepo repository.StateRepository) StateUsecase {
	return &stateUsecase{
		stateRepo: stateRepo,
	}
}

// Create создает новое состояние
func (u *stateUsecase) Create(s *state.State) error {
	logger.Info("Creating state in usecase with name: " + s.Name)

	if err := u.stateRepo.Create(s); err != nil {
		logger.Error("Failed to create state in repository", err)
		return err
	}

	logger.Info("State created successfully with ID: " + fmt.Sprintf("%d", s.ID))
	return nil
}

// GetByID получает состояние по ID
func (u *stateUsecase) GetByID(id uint) (*state.State, error) {
	return u.stateRepo.GetByID(id)
}

// GetAll получает список всех состояний
func (u *stateUsecase) GetAll() ([]*state.State, error) {
	return u.stateRepo.GetAll()
}

// Update обновляет существующее состояние
func (u *stateUsecase) Update(s *state.State) error {
	return u.stateRepo.Update(s)
}

// Delete удаляет состояние по ID
func (u *stateUsecase) Delete(id uint) error {
	return u.stateRepo.Delete(id)
}
