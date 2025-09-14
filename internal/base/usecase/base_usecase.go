package usecase

import (
	"github.com/Flood-project/backend-flood/internal/base"
	"github.com/Flood-project/backend-flood/internal/base/repository"
)

type BaseUseCase interface {
	Create(base *base.Base) error
	Fetch() ([]base.Base, error)
	Delete(id int32) error
}

type baseUseCase struct {
	baseRepository repository.BaseManagement
}

func NewBaseUseCase(baseRepository repository.BaseManagement) BaseUseCase {
	return &baseUseCase{
		baseRepository: baseRepository,
	}
}

func (baseUseCase *baseUseCase) Create(base *base.Base) error {
	err := baseUseCase.baseRepository.Create(base)
	if err != nil {
		return err
	}

	return nil
}

func (baseUseCase *baseUseCase) Fetch() ([]base.Base, error) {
	bases, err := baseUseCase.baseRepository.Fetch()
	if err != nil {
		return nil, err
	}

	return bases, nil
}

func (baseUseCase *baseUseCase) Delete(id int32) error {
	err := baseUseCase.baseRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}