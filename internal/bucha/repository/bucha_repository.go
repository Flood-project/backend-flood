package repository

import (
	"context"
	"log"

	"github.com/Flood-project/backend-flood/internal/bucha"
	"github.com/booscaaa/go-paginate/v3/paginate"
	"github.com/jmoiron/sqlx"
)

type BuchaManager interface {
	Create(bucha *bucha.Bucha) error
	Fetch() ([]bucha.Bucha, error)
	Delete(id int32) error
	GetWithParams(ctx context.Context, params *paginate.PaginationParams) ([]bucha.Bucha, int, error)
}

type buchaManager struct {
	DB *sqlx.DB
} 

func NewBuchaManager(db *sqlx.DB) BuchaManager{
	return &buchaManager{
		DB: db,
	}
}

func (buchaManager *buchaManager) Create(bucha *bucha.Bucha) error {
	query := `INSERT INTO buchas (tipobucha) VALUES ($1) RETURNING id`

	err := buchaManager.DB.QueryRow(query, bucha.TipoBucha).Scan(&bucha.ID)
	if err != nil {
		return err
	}

	return nil
}

func (buchaManager *buchaManager) Fetch() ([]bucha.Bucha, error) {
	query := `SELECT id, tipobucha FROM buchas`

	var buchas []bucha.Bucha
	err := buchaManager.DB.Select(&buchas, query)
	if err != nil {
		return nil, err
	}

	return buchas, nil
}

func (buchaManager *buchaManager) Delete(id int32) error {
	query := `DELETE FROM buchas WHERE id =$1`

	err := buchaManager.DB.QueryRow(query, id)
	if err != nil {
		return nil
	}
	return nil
}

func (buchaManager *buchaManager) GetWithParams(ctx context.Context, params *paginate.PaginationParams) ([]bucha.Bucha, int, error) {
	query, args, err := paginate.NewBuilder().
	Table("buchas").
	Model(&bucha.Bucha{}).
	FromStruct(params).
	BuildSQL()

	log.Println(query, args)

	if err != nil {
		return nil, 0, err
	}

	var buchasWithParams []bucha.Bucha

	err = buchaManager.DB.SelectContext(ctx, &buchasWithParams, query, args...)
	if err != nil {
		return nil, 0, err
	}

	total := len(buchasWithParams)

	return buchasWithParams, total, nil
}