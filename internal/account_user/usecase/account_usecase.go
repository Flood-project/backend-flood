package usecase

import (
	"github.com/Flood-project/backend-flood/internal/account_user"
	"github.com/Flood-project/backend-flood/internal/account_user/repository"
)

type AccountUseCase interface {
	Create(account account_user.Account) error
	Fetch() ([]account_user.Account, error)
}

type accountUseCase struct {
	saleRepository repository.AccountRepository
}

func (a *accountUseCase) Create(account *account_user.Account) error {
	err := a.saleRepository.Create(account)
	if err != nil {
		return err
	}

	return nil
}

func (a *accountUseCase) Fetch() ([]account_user.Account, error) {
	accounts, err := a.saleRepository.Fetch()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}