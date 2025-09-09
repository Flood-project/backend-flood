package repository

import (
	"github.com/Flood-project/backend-flood/internal/product"
	"github.com/jmoiron/sqlx"
)

type ProductManager interface {
	Create(product *product.Produt) error
	Fetch() ([]product.Produt, error)
	// GetByID(id int32) (*product.Produt, error)
	// Update(id int32, product *product.Produt) (*product.Produt, error)
	// Delete(id int32) error
}

type productManager struct {
	DB *sqlx.DB
}

func NewProductManager(db *sqlx.DB) ProductManager {
	return &productManager{
		DB: db,
	}
}

func (productManager *productManager) Create(product *product.Produt) error {
	query := `INSERT INTO products (
		name, description, id_bucha, id_acionamento, id_base, capacidade, valor
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`

	err := productManager.DB.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Id_bucha,
		product.Id_acionamento,
		product.Id_base,
		product.Capacity,
		product.Value,
	).Scan(&product.Id)
	if err != nil {
		return err
	}

	return nil
}

func (productManager *productManager) Fetch() ([]product.Produt, error) {
	query := `SELECT 
		id, name, description, id_bucha, id_acionamento, id_base, capacidade, valor
		FROM products`
	
	var products []product.Produt

	err := productManager.DB.Select(
		&products,
		query,
	)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (productManager *productManager) GetByID(id int32) (*product.Produt, error) {
	var product product.Produt
	query := `SELECT id, name, description, id_bucha, id_acionamento, id_base, capacidade, valor FROM products WHERE id = $1`

	err := productManager.DB.QueryRow(
		query,
		id,
	).Scan(
		&product.Id, 
		&product.Name,
		&product.Description,
		&product.Id_bucha,
		&product.Id_acionamento,
		&product.Id_base,
		&product.Capacity,
		&product.Value,
	)

	if err != nil {
		return nil, err
	}

	return &product, err
}