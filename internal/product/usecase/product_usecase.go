package usecase

import (
	"context"

	"github.com/Flood-project/backend-flood/config"
	"github.com/Flood-project/backend-flood/internal/object_store"
	"github.com/Flood-project/backend-flood/internal/product"
	"github.com/Flood-project/backend-flood/internal/product/repository"
	"github.com/booscaaa/go-paginate/v3/paginate"
)

type ProductUseCase interface {
	Create(product *product.Produt) (*product.Produt, error)
	Fetch() ([]product.Produt, error)
	FetchWithComponents() ([]product.ProductWithComponents, error)
	GetByID(id int32) (*product.ProductWithComponents, error)
	Update(id int32, product *product.Produt) error
	Delete(id int32) error
	WithParams(ctx context.Context, params *paginate.PaginationParams) ([]product.ProductWithComponents, config.PageData, error)
	GetProductByIdWithImage(id int32) ([]object_store.FileData, error)
}

type productUseCase struct {
	productRepository repository.ProductManager
}

func NewProductUseCase(productRepository *repository.ProductManager) ProductUseCase {
	return &productUseCase{
		productRepository: *productRepository,
	}
}

func (productUseCase *productUseCase) Create(product *product.Produt) (*product.Produt, error) {
	product, err := productUseCase.productRepository.Create(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (productUseCase *productUseCase) Fetch() ([]product.Produt, error) {
	products, err := productUseCase.productRepository.Fetch()
	if err != nil {
		return nil, err
	}

	return products, err
}

func (productUseCase *productUseCase) FetchWithComponents() ([]product.ProductWithComponents, error) {
	productsWithComponents, err := productUseCase.productRepository.FetchWithComponents()
	if err != nil {
		return nil, err
	}

	return productsWithComponents, err
}

func (productUseCase *productUseCase) GetByID(id int32) (*product.ProductWithComponents, error) {
	product, err := productUseCase.productRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (productUseCase *productUseCase) Update(id int32, product *product.Produt) error {
	err := productUseCase.productRepository.Update(id, product)
	if err != nil {
		return err
	}

	return nil
}

func (productUseCase *productUseCase) Delete(id int32) error {
	err := productUseCase.productRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (productUseCase *productUseCase) WithParams(ctx context.Context, params *paginate.PaginationParams) ([]product.ProductWithComponents, config.PageData, error) {
	rows, total, err := productUseCase.productRepository.WithParams(ctx, params)
	if err != nil {
		return nil, config.PageData{}, err
	}

	return rows, config.PageData{
		Total: int64(total),
		Page:  int64(params.Page),
		Limit: int64(params.ItemsPerPage),
	}, nil
}

func (productUseCase *productUseCase) GetProductByIdWithImage(productID int32) ([]object_store.FileData, error) {
	image, err := productUseCase.productRepository.GetProductByIdWithImage(productID)
	if err != nil {
		return nil, err
	}

	return image, nil
}
