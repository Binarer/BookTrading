package usecase

import (
	"booktrading/internal/domain/state"
	"booktrading/internal/repository"
)

type StateUsecase interface {
	CreateState(dto *state.CreateStateDTO) (*state.State, error)
	GetStateByID(id int64) (*state.State, error)
	GetAllStates() ([]*state.State, error)
	UpdateState(id int64, dto *state.UpdateStateDTO) (*state.State, error)
	DeleteState(id int64) error
}

type stateUsecase struct {
	stateRepo repository.StateRepository
}

func NewStateUsecase(stateRepo repository.StateRepository) StateUsecase {
	return &stateUsecase{
		stateRepo: stateRepo,
	}
}

func (u *stateUsecase) CreateState(dto *state.CreateStateDTO) (*state.State, error) {
	s := &state.State{
		Name: dto.Name,
	}

	if err := u.stateRepo.Create(s); err != nil {
		return nil, err
	}

	return s, nil
}

func (u *stateUsecase) GetStateByID(id int64) (*state.State, error) {
	return u.stateRepo.GetByID(id)
}

func (u *stateUsecase) GetAllStates() ([]*state.State, error) {
	return u.stateRepo.GetAll()
}

func (u *stateUsecase) UpdateState(id int64, dto *state.UpdateStateDTO) (*state.State, error) {
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

func (u *stateUsecase) DeleteState(id int64) error {
	return u.stateRepo.Delete(id)
} 