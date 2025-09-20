package usecase

import (
	"github.com/Flood-project/backend-flood/internal/bucha"
	"github.com/Flood-project/backend-flood/internal/bucha/repository"
)

type BuchaUseCase interface {
	Create(bucha *bucha.Bucha) error
	Fetch() ([]bucha.Bucha, error)
	Delete(id int32) error
}

type buchaUseCase struct {
	buchaRepository repository.BuchaManager
}

func NewBuchaUseCase(buchaRepository repository.BuchaManager) BuchaUseCase {
	return &buchaUseCase{
		buchaRepository: buchaRepository,
	}
}

func (usecase *buchaUseCase) Create(bucha *bucha.Bucha) error {
	err := usecase.buchaRepository.Create(bucha)
	if err != nil {
		return nil
	}
	
	return nil
}

func (usecase *buchaUseCase) Fetch() ([]bucha.Bucha, error) {
	buchas, err := usecase.buchaRepository.Fetch()
	if err != nil {
		return nil, err
	}

	return buchas, nil
}

func (usecase *buchaUseCase) Delete(id int32) error {
	err := usecase.buchaRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}