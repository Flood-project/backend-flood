package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/Flood-project/backend-flood/internal/object_store"
	"github.com/Flood-project/backend-flood/internal/object_store/repository"
	"github.com/Flood-project/backend-flood/internal/product"
	"github.com/booscaaa/go-paginate/v3/paginate"
	"github.com/jmoiron/sqlx"
)

type ProductManager interface {
	Create(product *product.Produt) (*product.Produt, error)
	Fetch() ([]product.Produt, error)
	FetchWithComponents() ([]product.ProductWithComponents, error)
	GetByID(id int32) (*product.ProductWithComponents, error)
	Update(id int32, product *product.Produt) error
	Delete(id int32) error
	WithParams(ctx context.Context, params *paginate.PaginationParams) ([]product.ProductWithComponents, int, error)
	GetProductByIdWithImage(id int32) ([]object_store.FileData, error)
}

type productManager struct {
	DB                 *sqlx.DB
	ObjectStoreManager repository.ObjectStoreManager
}

func NewProductManager(db *sqlx.DB, objectStore repository.ObjectStoreManager) ProductManager {
	return &productManager{
		DB:                 db,
		ObjectStoreManager: objectStore,
	}
}

func (productManager *productManager) Create(product *product.Produt) (*product.Produt, error) {
	query := `INSERT INTO products (
		codigo, description, capacidade_estatica, capacidade_trabalho, reducao, altura_bucha, curso, id_bucha, id_acionamento, id_base, ativo
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
		product.Id_base,
		product.Ativo,
	).Scan(&product.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return product, nil
}

func (productManager *productManager) Fetch() ([]product.Produt, error) {
	query := `SELECT 
		id, codigo, description, capacidade_estatica, capacidade_trabalho, reducao, altura_bucha, curso, id_bucha, id_acionamento, id_base, ativo
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

func (productManager *productManager) GetProductByIdWithImage(productID int32) ([]object_store.FileData, error) {
	query := `
        SELECT 
            f.id,
            f.product_id,
            f.file_name,
            f.storage_key,
            f.size,
            f.content_type,
            'http://localhost:8080/files/images/' || f.storage_key as url
        FROM files f
        WHERE f.product_id = $1
    `

	var files []object_store.FileData
	err := productManager.DB.Select(&files, query, productID)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (productManager *productManager) FetchWithComponents() ([]product.ProductWithComponents, error) {
	query := `SELECT 
				p.id,
				p.codigo,
				p.description,c
				p.capacidade_estatica,
				p.capacidade_trabalho,
				p.reducao,
				p.altura_bucha,
				p.curso,
				p.ativo,
				b.id AS id_bucha,
				b.tipobucha,
				a.id AS id_acionamento,
				a.tipoacionamento,
				base.id AS id_base,
				base.tipoBase
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

	for i := range productsWithComponents {
		images, err := productManager.GetProductByIdWithImage(int32(productsWithComponents[i].Id))
		if err != nil {
			log.Println("Erro ao buscar imagens do produto... ", err)
			continue
		}
		productsWithComponents[i].Images = images
	}

	return productsWithComponents, nil
}

func (productManager *productManager) GetByID(id int32) (*product.ProductWithComponents, error) {
	var product product.ProductWithComponents
	var images []object_store.FileData
	first := true
	// query := `
	// 		SELECT id,
	// 		codigo,
	// 		description,
	// 		capacidade_estatica,
	// 		capacidade_trabalho,
	// 		reducao,
	// 		altura_bucha,
	// 		curso,
	// 		ativo,
	// 		id_bucha,
	// 		id_acionamento,
	// 		id_base
	// 			FROM products WHERE id = $1
	// 	`
	query := `
		SELECT 
		  f.id,
		  f.product_id,
		  p.codigo,
			p.description, 
			p.capacidade_estatica, 
			p.capacidade_trabalho, 
			p.reducao, 
			p.altura_bucha, 
			p.curso, 
			p.ativo, 
			p.id_bucha, 
			p.id_acionamento, 
			p.id_base, 
		  f.file_name,
		  'http://localhost:8080/files/images/' || f.storage_key as url,
		  f.size,
		  f.content_type
		FROM files f
		INNER JOIN products p
		  ON p.id = f.product_id WHERE p.id = $1
	`

	rows, err := productManager.DB.Query(
		query,
		id,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var file object_store.FileData

		err := rows.Scan(
			&file.ID,
			&file.ProductID,
			&product.Codigo,
			&product.Description,
			&product.CapacidadeEstatica,
			&product.CapacidadeTrabalho,
			&product.Reducao,
			&product.AlturaBucha,
			&product.Curso,
			&product.Ativo,
			&product.IdBucha,
			&product.IdAcionamento,
			&product.IdBase,
			&file.FileName,
			&file.URL,
			&file.Size,
			&file.ContentType,
		)
		if err != nil {
			return nil, err
		}

		if first {
			product.Id = file.ProductID
			first = false
		}

		images = append(images, file)
	}

	product.Images = images

	return &product, nil
}

func (productManager *productManager) Update(id int32, product *product.Produt) error {
	query := `UPDATE products SET codigo=$1, description=$2, capacidade_estatica=$3, capacidade_trabalho=$4, reducao=$5, altura_bucha=$6, curso=$7, id_bucha=$8, id_acionamento=$9, id_base=$10, ativo=$11 WHERE id=$12`

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
		product.Id_base,
		product.Ativo,
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

	result, err := productManager.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao executar delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("produto com id %d não encontrado", id)
	}

	return nil
}

func (produtctManager *productManager) WithParams(ctx context.Context, params *paginate.PaginationParams) ([]product.ProductWithComponents, int, error) {
	query, args, err := paginate.NewBuilder().
		Table("products p").
		Model(&product.ProductWithComponents{}).
		Select("p.*", "b.tipobucha", "a.tipoacionamento", "bs.tipobase").
		LeftJoin("buchas b", "p.id_bucha = b.id").
		LeftJoin("acionamentos a", "p.id_acionamento = a.id").
		LeftJoin("bases bs", "p.id_base = bs.id").
		FromStruct(params).
		BuildSQL()

	if err != nil {
		return nil, 0, err
	}

	var products []product.ProductWithComponents

	err = produtctManager.DB.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, 0, err
	}

	for i := range products {
		images, err := produtctManager.GetProductByIdWithImage(int32(products[i].Id))
		if err != nil {
			log.Println("Erro ao buscar imagens do produto... ", err)
			continue
		}
		products[i].Images = images
	}

	total := len(products)

	return products, total, nil
}
