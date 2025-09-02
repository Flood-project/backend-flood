package repository

import (
	"github.com/Flood-project/backend-flood/internal/token"
	"github.com/jmoiron/sqlx"
)

type TokenRepository interface {
	Create(token *token.Token) error
	Fetch() ([]token.Token, error)
}

type tokenRepository struct {
	DB *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) TokenRepository {
	return &tokenRepository{
		DB: db,
	}
}

func (tokenRepository *tokenRepository) Create(token *token.Token) error {
	query := `INSERT INTO auth 
	(token, created, expiration, id_account) 
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	err := tokenRepository.DB.QueryRow(
		query,
		&token.RowToken,
		&token.Created,
		&token.Expiration,
		&token.IdAccount,
	).Scan(
		&token.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (tokenRepository *tokenRepository) Fetch() ([]token.Token, error) {
	query := `SELECT id, token, created, expiration, id_account FROM auth`

	var tokens []token.Token
	err := tokenRepository.DB.Select(&tokens, query)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}