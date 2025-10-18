package repository

import (
	"context"
	"fmt"
	"github.com/Flood-project/backend-flood/internal/product"
	"github.com/jmoiron/sqlx"
)

type ProductManager interface {
	Create(product *product.Produt) error
	Fetch() ([]product.Produt, error)
	FetchWithComponents() ([]product.ProductWithComponents, error)
	GetByID(id int32) (*product.Produt, error)
	Update(id int32, product *product.Produt) error
	Delete(id int32) error
	WithParams(ctx context.Context, query string, args ...interface{}) ([]product.ProductWithComponents, error)
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
		codigo, description, capacidade_estatica, capacidade_trabalho, reducao, altura_bucha, curso, id_bucha, id_acionamento, id_base
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`

	err := productManager.DB.QueryRow(
		query,
		product.Codigo,
		product.Description,
		product.CapacidadeEstatica,
		product.CapacidadeTrabalho,
		product.Reducao,
		product.AlturaBucha,
		product.Curso,
		product.Id_bucha,
		product.Id_acionamento,
		//product.Id_base,
	).Scan(&product.Id)
	if err != nil {
		return err
	}

	return nil
}

func (productManager *productManager) Fetch() ([]product.Produt, error) {
	query := `SELECT 
		id, codigo, description, capacidade_estatica, capacidade_trabalho, reducao, altura_bucha, curso, id_bucha, id_acionamento, id_base
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

func (productManager *productManager) FetchWithComponents() ([]product.ProductWithComponents, error) {
	query := `SELECT 
				p.id,
				p.codigo,
				p.description,
				p.capacidade_estatica,
				p.capacidade_trabalho,
				p.reducao,
				p.altura_bucha,
				p.curso,
				b.id AS id_bucha,
				b.tipobucha AS tipo_bucha,
				a.id AS id_acionamento,
				a.tipoacionamento AS tipo_do_acionamento,
				base.id AS id_base,
				base.tipoBase AS tipo_base
			  FROM products p
			  INNER JOIN buchas b
				ON b.id = p.id_bucha
			  INNER JOIN acionamentos a
				ON a.id = p.id_acionamento
			  INNER JOIN bases base
			  	ON base.id = p.id_base`
	
	var productsWithComponents []product.ProductWithComponents

	err := productManager.DB.Select(
		&productsWithComponents,
		query,
	)
	if err != nil {
		return nil, err
	}

	return productsWithComponents, nil
}

func (productManager *productManager) GetByID(id int32) (*product.Produt, error) {
	var product product.Produt
	query := `SELECT id, codigo, description, capacidade_estatica, capacidade_trabalho, reducao, altura_bucha, curso, id_bucha, id_acionamento, id_base FROM products WHERE id = $1`

	err := productManager.DB.QueryRow(
		query,
		id,
	).Scan(
		&product.Id, 
		&product.Codigo,
		&product.Description,
		&product.CapacidadeEstatica,
		&product.CapacidadeTrabalho,
		&product.Reducao,
		&product.AlturaBucha,
		&product.Curso,
		&product.Id_bucha,
		&product.Id_acionamento,
		&product.Id_base,
	)

	if err != nil {
		return nil, err
	}

	return &product, err
}

func (productManager *productManager) Update(id int32, product *product.Produt) error {
	query := `UPDATE products SET codigo=$1, description=$2, capacidade_estatica=$3, capacidade_trabalho=$4, reducao=$5, altura_bucha=$6, curso=$7, id_bucha=$8, id_acionamento=$9, id_base=$10 WHERE id=$11`

	res, err := productManager.DB.Exec(
		query,
		product.Codigo,
		product.Description,
		product.CapacidadeEstatica,
		product.CapacidadeTrabalho,
		product.Reducao,
		product.AlturaBucha,
		product.Curso,
		product.Id_bucha,
		product.Id_acionamento,
		//product.Id_base,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Error: Nenhum produto foi modificado.")
	}

	return nil
}

func (productManager *productManager) Delete(id int32) error {
	query := `DELETE FROM products WHERE id=$1`

	err := productManager.DB.QueryRow(
		query,
		id,
	)
	if err != nil {
		return nil
	}
	return nil
}

func (produtctManager *productManager) WithParams(ctx context.Context, query string, args ...interface{}) ([]product.ProductWithComponents, error) {
	var products []product.ProductWithComponents

	err := produtctManager.DB.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, err
	}
	return products, nil
}