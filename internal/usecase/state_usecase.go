package usecase

import (
	"booktrading/internal/domain/repository"
	"booktrading/internal/domain/state"
	"booktrading/internal/pkg/logger"
	"fmt"
)

// StateUseCase определяет интерфейс для работы с состояниями книг
type StateUseCase interface {
	Create(s *state.State) error
	GetByID(id uint) (*state.State, error)
	GetAll() ([]*state.State, error)
	Update(s *state.State) error
	Delete(id uint) error
}

// stateUseCase реализует интерфейс StateUseCase
type stateUseCase struct {
	stateRepo repository.StateRepository
}

// NewStateUseCase создает новый экземпляр stateUseCase
func NewStateUseCase(stateRepo repository.StateRepository) StateUseCase {
	return &stateUseCase{
		stateRepo: stateRepo,
	}
}

// Create создает новое состояние
func (u *stateUseCase) Create(s *state.State) error {
	logger.Info("Creating state in usecase with name: " + s.Name)

	if err := u.stateRepo.Create(s); err != nil {
		logger.Error("Failed to create state in repository", err)
		return err
	}

	logger.Info("State created successfully with ID: " + fmt.Sprintf("%d", s.ID))
	return nil
}

// GetByID получает состояние по ID
func (u *stateUseCase) GetByID(id uint) (*state.State, error) {
	return u.stateRepo.GetByID(id)
}

// GetAll получает список всех состояний
func (u *stateUseCase) GetAll() ([]*state.State, error) {
	return u.stateRepo.GetAll()
}

// Update обновляет существующее состояние
func (u *stateUseCase) Update(s *state.State) error {
	return u.stateRepo.Update(s)
}

// Delete удаляет состояние по ID
func (u *stateUseCase) Delete(id uint) error {
	return u.stateRepo.Delete(id)
}
