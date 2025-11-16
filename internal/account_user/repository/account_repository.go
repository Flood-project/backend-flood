package repository

import (
	"fmt"

	"github.com/Flood-project/backend-flood/internal/account_user"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	Create(account *account_user.Account) error
	Fetch() (accounts []account_user.Account, err error)
	FetchWithUserGroup() ([]account_user.AccountWithUserGroup, error)
	GetByID(id int32) (*account_user.Account, error)
	GetByEmail(email string) (*account_user.Account, error)
	GetUserGroup() ([]account_user.AccountGroupName, error)
	UpdateUser(id int32, account *account_user.Account) error
	DeleteUser(id int32) error
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

func (ur *accountRepository) FetchWithUserGroup() ([]account_user.AccountWithUserGroup, error) {
	query := `
		SELECT
			a.id, 
			a.name,
			a.email,
			a.password_hash,
			ug.id AS id_user_group,
			ug.group_name AS group_name
		FROM account a
		INNER JOIN user_group ug
			ON ug.id = a.id_user_group
	`

	var accountsWithUserGroup []account_user.AccountWithUserGroup

	err := ur.DB.Select(&accountsWithUserGroup, query)
	if err != nil {
		return nil, err
	}

	return accountsWithUserGroup, nil
}

func (ur *accountRepository) GetByID(id int32) (*account_user.Account, error) {
	var account account_user.Account

	query := `SELECT id, name, email, password_hash, id_user_group FROM account WHERE id = $1`

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

	query := `SELECT id, email, password_hash, id_user_group FROM account WHERE email = $1`

	err := accountRepository.DB.QueryRow(
		query,
		email,
	).Scan(&account.Id_account, &account.Email, &account.Password_hash, &account.IdUserGroup)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (repository *accountRepository) GetUserGroup() ([]account_user.AccountGroupName, error) {
	var accountGorupNames []account_user.AccountGroupName
	query := `SELECT id, group_name FROM user_group`

	err := repository.DB.Select(&accountGorupNames, query)
	if err != nil {
		return nil, err
	}

	return accountGorupNames, nil
}

func (repository *accountRepository) UpdateUser(id int32, account *account_user.Account) error {
	query := `UPDATE account SET name=$1, id_user_group=$2 WHERE id=$3`

	res, err := repository.DB.Exec(
		query,
		account.Name,
		account.IdUserGroup,
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
		return fmt.Errorf("Error: Nenhuma conta foi modificada.")
	}

	return nil
}

func (repository *accountRepository) DeleteUser(id int32) error {
	query := `DELETE FROM account WHERE id=$1`

	err := repository.DB.QueryRow(
		query,
		id,
	)
	if err != nil {
		return nil
	}
	return nil
}