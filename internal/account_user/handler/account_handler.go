package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Flood-project/backend-flood/internal/account_user"
	"github.com/Flood-project/backend-flood/internal/account_user/usecase"
	"github.com/Flood-project/backend-flood/internal/token"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type AccountHandler struct {
	accountUsecase usecase.AccountUseCase
	token          token.TokenManager
}

func NewAccountHandler(accountUsecase usecase.AccountUseCase, token token.TokenManager) *AccountHandler {
	return &AccountHandler{
		accountUsecase: accountUsecase,
		token:          token,
	}
}

func (handler *AccountHandler) Create(response http.ResponseWriter, request *http.Request) {
	var account account_user.Account

	err := json.NewDecoder(request.Body).Decode(&account)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(account.Password_hash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	hashedAccount := account_user.Account{
		Name:          account.Name,
		Email:         account.Email,
		Password_hash: string(hashedPass),
		IdUserGroup:   account.IdUserGroup,
	}

	err = handler.accountUsecase.Create(&hashedAccount)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	hashedAccount.Password_hash = ""
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(hashedAccount)
	if err != nil {
		http.Error(response, "Erro ao criar nova conta", http.StatusBadRequest)
	}
}

func (handler *AccountHandler) Fetch(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	tokenHeader := request.Header.Get("Authorization")
	if tokenHeader == "" {
		http.Error(response, "Usuário não autorizado.", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(tokenHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(response, "Formato do token inválido.", http.StatusUnauthorized)
		return
	}
	tokenString := parts[1]

	_, err := handler.token.ValidateToken(tokenString)
	if err != nil {
		http.Error(response, "Token inválido.", http.StatusUnauthorized)
		return
	}

	accounts, err := handler.accountUsecase.Fetch()
	if err != nil {
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(response).Encode(&accounts)
	if err != nil {
		http.Error(response, "Erro nos dados json", http.StatusBadRequest)
		return
	}
}

func (handler *AccountHandler) FetchWithUserGroup(response http.ResponseWriter, request *http.Request) {
	products, err := handler.accountUsecase.FetchWithUserGroup()
	if err != nil {
		http.Error(response, "Erro ao listar contas com grupo de usuário", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err = json.NewEncoder(response).Encode(&products)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *AccountHandler) GetByID(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, "Erro ao listar por id", http.StatusBadRequest)
		return
	}

	account, err := handler.accountUsecase.GetByID(int32(id))
	if err != nil {
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(account)
}

func (handler *AccountHandler) GetUserGroup(response http.ResponseWriter, request *http.Request) {
	userGroupName, err := handler.accountUsecase.GetUserGroup()
	if err != nil {
		log.Println(err)
		http.Error(response, "Erro ao listar nome do grupo de usuário. ", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err = json.NewEncoder(response).Encode(&userGroupName)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *AccountHandler) UpdateAccount(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, "Conta não encontrada. ", http.StatusBadRequest)
		return
	}

	var account account_user.Account
	err = json.NewDecoder(request.Body).Decode(&account)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.accountUsecase.UpdateUser(int32(id), &account)
	if err != nil {
		http.Error(response, "Erro ao atualizar um usuário. ", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&account)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *AccountHandler) DeleteAccount(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, "Conta não encontrada. ", http.StatusBadRequest)
		return
	}

	err = handler.accountUsecase.DeleteUser(int32(id))
	if err != nil {
		http.Error(response, "Erro ao deletar um usuário. ", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}
