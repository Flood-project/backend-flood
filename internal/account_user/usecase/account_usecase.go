package usecase

import (
	"github.com/Flood-project/backend-flood/internal/account_user"
	"github.com/Flood-project/backend-flood/internal/account_user/repository"
)

type AccountUseCase interface {
	Create(account *account_user.Account) error
	Fetch() ([]account_user.Account, error)
	FetchWithUserGroup() ([]account_user.AccountWithUserGroup, error)
	GetByID(id int32) (*account_user.Account, error)
	GetByEmail(email string) (*account_user.Account, error)
}

type accountUseCase struct {
	accountRepository repository.AccountRepository
}

func (a *accountUseCase) Create(account *account_user.Account) error {
	err := a.accountRepository.Create(account)
	if err != nil {
		return err
	}

	return nil
}

func (a *accountUseCase) Fetch() ([]account_user.Account, error) {
	accounts, err := a.accountRepository.Fetch()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a *accountUseCase) FetchWithUserGroup() ([]account_user.AccountWithUserGroup, error) {
	accountsWithUserGroup, err := a.accountRepository.FetchWithUserGroup()
	if err != nil {
		return nil, err
	}

	return accountsWithUserGroup, nil
}

func (a *accountUseCase) GetByID(id int32) (*account_user.Account, error) {
	account, err := a.accountRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return account, nil	
}

func (a *accountUseCase) GetByEmail(email string) (*account_user.Account, error) {
	account, err := a.accountRepository.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return account, nil
}