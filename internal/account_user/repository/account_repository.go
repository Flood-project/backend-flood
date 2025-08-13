package repository

import (
	"github.com/Flood-project/backend-flood/internal/account_user"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	Create(account *account_user.Account) error
	Fetch() (accounts []account_user.Account, err error)
}

type accountRepository struct {
	DB *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &accountRepository{
		DB: db,
	}  
}

func (ur *accountRepository) Create(account *account_user.Account) error {
	query := `INSERT INTO account (name, email, password_hash)
		VALUES($1, $2, $3, $4)
		RETURNING id_account`

	err := ur.DB.QueryRow(
		query,
		account.Name,
		account.Email,
		account.Password_hash,
	).Scan(&account.Id_account)

	if err != nil {
		return err
	}

	return nil
}

func (ur *accountRepository) Fetch() (accounts []account_user.Account, err error) {
	query := `SELECT (id_account, name, email, password_hash) FROM account`

	err = ur.DB.Select(&accounts, query)
	if err != nil {
		return
	}

	return
}