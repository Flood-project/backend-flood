package usecase

import (
	"context"

	"github.com/Flood-project/backend-flood/config"
	"github.com/Flood-project/backend-flood/internal/bucha"
	"github.com/Flood-project/backend-flood/internal/bucha/repository"
	"github.com/booscaaa/go-paginate/v3/paginate"
)

type BuchaUseCase interface {
	Create(bucha *bucha.Bucha) error
	Fetch() ([]bucha.Bucha, error)
	Delete(id int32) error
	GetWithParams(ctx context.Context, params *paginate.PaginationParams) ([]bucha.Bucha, config.PageData, error)
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

func (usecase *buchaUseCase) GetWithParams(ctx context.Context, params *paginate.PaginationParams) ([]bucha.Bucha, config.PageData, error) {
	buchasWithParams, total, err := usecase.buchaRepository.GetWithParams(ctx, params)
	if err != nil {
		return nil, config.PageData{}, err
	}

	return buchasWithParams, config.PageData{
		Total: int64(total),
		Page: int64(params.Page),
		Limit: int64(params.ItemsPerPage),
	}, nil
}