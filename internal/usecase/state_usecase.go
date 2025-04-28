package usecase

import (
	"booktrading/internal/domain/state"
	"booktrading/internal/repository"
)

// StateUsecase определяет интерфейс для работы с состояниями книг
type StateUsecase interface {
	CreateState(dto *state.CreateStateDTO) (*state.State, error)
	GetStateByID(id uint) (*state.State, error)
	GetAllStates() ([]*state.State, error)
	UpdateState(id uint, dto *state.UpdateStateDTO) (*state.State, error)
	DeleteState(id uint) error
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

// CreateState создает новое состояние
func (u *stateUsecase) CreateState(dto *state.CreateStateDTO) (*state.State, error) {
	s := &state.State{
		Name: dto.Name,
	}

	if err := u.stateRepo.Create(s); err != nil {
		return nil, err
	}

	return s, nil
}

// GetStateByID получает состояние по ID
func (u *stateUsecase) GetStateByID(id uint) (*state.State, error) {
	return u.stateRepo.GetByID(id)
}

// GetAllStates получает список всех состояний
func (u *stateUsecase) GetAllStates() ([]*state.State, error) {
	return u.stateRepo.GetAll()
}

// UpdateState обновляет существующее состояние
func (u *stateUsecase) UpdateState(id uint, dto *state.UpdateStateDTO) (*state.State, error) {
	s, err := u.stateRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	s.Name = dto.Name

	if err := u.stateRepo.Update(s); err != nil {
		return nil, err
	}

	return s, nil
}

// DeleteState удаляет состояние по ID
func (u *stateUsecase) DeleteState(id uint) error {
	return u.stateRepo.Delete(id)
} 