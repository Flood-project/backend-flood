package repository

import (
	"fmt"

	"github.com/Flood-project/backend-flood/internal/acionameto"
	"github.com/jmoiron/sqlx"
)

type AcionamentoManagement interface {
	Create(acionamento *acionametos.Acionamento) error
	Fetch() ([]acionametos.Acionamento, error)
	Delete(id int32) error
	UpdateAcionamento(id int32, acionamento *acionametos.Acionamento) error
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

func (acionamentoManagement *acionamentoManagement) UpdateAcionamento(id int32, acionamento *acionametos.Acionamento) error {
	query := `UPDATE acionamentos SET tipoacionamento=$1 WHERE id=$2`

	res, err := acionamentoManagement.DB.Exec(
		query,
		acionamento.TipoAcionamento,
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
		return fmt.Errorf("Error: Nenhum acionamento foi modificado.")
	}

	return nil
}