package repository

import (
	"github.com/Flood-project/backend-flood/internal/base"
	"github.com/jmoiron/sqlx"
)

type BaseManagement interface {
	Create(base *base.Base) error
	Fetch() ([]base.Base, error)
	Delete(id int32) error
}

type baseManagement struct {
	DB *sqlx.DB
}

func NewBaseManagement(db *sqlx.DB) BaseManagement {
	return &baseManagement{
		DB: db,
	}
}

func (management *baseManagement) Create(base *base.Base) error {
	query := `INSERT INTO bases (tipobase) VALUES ($1) RETURNING id`

	err := management.DB.QueryRow(query, base.TipoBase).Scan(&base.ID)
	if err != nil {
		return err
	}

	return nil
}

func (management *baseManagement) Fetch() ([]base.Base, error) {
	query := `SELECT id, tipobase FROM bases`

	var bases []base.Base

	err := management.DB.Select(&bases, query)
	if err != nil {
		return nil, err
	}

	return bases, err
}

func (management *baseManagement) Delete(id int32) error {
	query := `DELETE FROM bases WHERE id = $1`

	err := management.DB.QueryRow(query, id)
	if err != nil {
		return nil
	}

	return nil
}

