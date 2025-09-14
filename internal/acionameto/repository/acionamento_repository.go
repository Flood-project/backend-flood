package repository

import (
	"github.com/Flood-project/backend-flood/internal/acionameto"
	"github.com/jmoiron/sqlx"
)

type AcionamentoManagement interface {
	Create(acionamento *acionametos.Acionamento) error
	Fetch() ([]acionametos.Acionamento, error)
	Delete(id int32) error
}

type acionamentoManagement struct {
	DB *sqlx.DB
}

func NewAcionamentoManagement(db *sqlx.DB) AcionamentoManagement {
	return &acionamentoManagement{
		DB: db,
	}
}

func (acionamentoManagement *acionamentoManagement) Create(acionamento *acionametos.Acionamento) error {
	query := `INSERT INTO acionamentos (tipoacionamento) VALUES ($1) RETURNING id`

	err := acionamentoManagement.DB.QueryRow(query, acionamento.TipoAcionamento).Scan(&acionamento.ID)
	if err != nil {
		return err
	}

	return nil
}

func (acionamentoManagement *acionamentoManagement) Fetch() ([]acionametos.Acionamento, error) {
	query := `SELECT id, tipoacionamento FROM acionamentos`

	var acionamentos []acionametos.Acionamento

	err := acionamentoManagement.DB.Select(&acionamentos, query)
	if err != nil {
		return nil, err
	}
	return acionamentos, nil
}

func (acionamentoManagement *acionamentoManagement) Delete(id int32) error {
	query := `DELETE FROM acionamentos WHERE id = $1`

	err := acionamentoManagement.DB.QueryRow(query, id)
	if err != nil {
		return nil
	}

	return nil
}