package handler

import (
	"encoding/json"
	"net/http"
	"github.com/Flood-project/backend-flood/internal/account_user"
	"github.com/Flood-project/backend-flood/internal/account_user/usecase"
	"golang.org/x/crypto/bcrypt"
)

type AccountHandler struct {
	accountUsecase usecase.AccountUseCase
}

func NewAccountHandler(accountUsecase usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{
		accountUsecase: accountUsecase,
	}
}

func (handler *AccountHandler) Create(response http.ResponseWriter, request *http.Request) {
	var account account_user.Account



	err := json.NewDecoder(request.Body).Decode(&account)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(account.Password_hash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	hashedAccount := account_user.Account{
		Name: account.Name,
		Email: account.Email,
		Password_hash: string(hashedPass),
	}

	err = handler.accountUsecase.Create(&hashedAccount)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	hashedAccount.Password_hash = ""
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(hashedAccount)
	if err != nil {
		http.Error(response, "Erro ao criar nova conta", http.StatusBadRequest)
	}
}

func (handler *AccountHandler) Fetch(response http.ResponseWriter, request *http.Request) {
	accounts, err := handler.accountUsecase.Fetch()
	if err != nil {
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(response).Encode(&accounts)
	if err != nil {
		http.Error(response, "Erro nos dados json", http.StatusBadRequest)
	}
}
