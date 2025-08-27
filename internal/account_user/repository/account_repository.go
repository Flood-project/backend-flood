package repository

import (
	"github.com/Flood-project/backend-flood/internal/account_user"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	Create(account *account_user.Account) error
	Fetch() (accounts []account_user.Account, err error)
	GetByID(id int32) (*account_user.Account, error)
	GetByEmail(email string) (*account_user.Account, error)
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
	query := `INSERT INTO account (name, email, password_hash, id_user_group)
		VALUES($1, $2, $3, $4)
		RETURNING id`

	err := ur.DB.QueryRow(
		query,
		account.Name,
		account.Email,
		account.Password_hash,
		account.IdUserGroup,
	).Scan(&account.Id_account)

	if err != nil {
		return err
	}

	return nil
}

func (ur *accountRepository) Fetch() (accounts []account_user.Account, err error) {
	query := `SELECT id, name, email, password_hash, id_user_group FROM account`

	err = ur.DB.Select(&accounts, query)
	if err != nil {
		return
	}

	return
}

func (ur *accountRepository) GetByID(id int32) (*account_user.Account, error) {
	var account account_user.Account

	query := `SELECT id, name, email, password_hash FROM account WHERE id = $1`

	err := ur.DB.QueryRow(
		query,
		id,
	).Scan(
		&account.Id_account,
		&account.Name,
		&account.Email,
		&account.Password_hash,
		&account.IdUserGroup,
	)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (accountRepository *accountRepository) GetByEmail(email string) (*account_user.Account, error) {
	var account account_user.Account

	query := `SELECT email, password_hash, id_user_group FROM account WHERE email = $1`

	err := accountRepository.DB.QueryRow(
		query,
		email,
	).Scan(&account.Email, &account.Password_hash, &account.IdUserGroup)
	if err != nil {
		return nil, err
	}

	return &account, nil
}