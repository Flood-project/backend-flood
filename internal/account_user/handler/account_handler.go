package handler

import "github.com/Flood-project/backend-flood/internal/account_user/usecase"

type AccountHandler struct {
	accountUsecase usecase.AccountUseCase
}

func NewAccountHandler(accountUsecase usecase.AccountUseCase) *AccountHandler{
	return &AccountHandler{
		accountUsecase: accountUsecase,
	}
}