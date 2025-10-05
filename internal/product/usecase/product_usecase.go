package usecase

import (
	"context"
	"github.com/Flood-project/backend-flood/internal/product"
	"github.com/Flood-project/backend-flood/internal/product/repository"
)

type ProductUseCase interface {
	Create(product *product.Produt) error
	Fetch() ([]product.Produt, error)
	GetByID(id int32) (*product.Produt, error)
	Update(id int32, product *product.Produt) error
	Delete(id int32) error
	WithParams(ctx context.Context, query string, args ...interface{}) ([]product.ProductWithComponents, error)

}

type productUseCase struct {
	productRepository repository.ProductManager
}

func NewProductUseCase(productRepository *repository.ProductManager) ProductUseCase{
	return &productUseCase{
		productRepository: *productRepository,
	}
}

func (productUseCase *productUseCase) Create(product *product.Produt) error {
	err := productUseCase.productRepository.Create(product)
	if err != nil {
		return nil
	}
	return nil
}

func (productUseCase *productUseCase) Fetch() ([]product.Produt, error) {
	products, err := productUseCase.productRepository.Fetch()
	if err != nil {
		return nil, err
	}

	return products, err
}

func (productUseCase *productUseCase) GetByID(id int32) (*product.Produt, error) {
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

func (productUseCase *productUseCase) WithParams(ctx context.Context, query string, args ...interface{}) ([]product.ProductWithComponents, error) {
	rows, err := productUseCase.productRepository.WithParams(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
