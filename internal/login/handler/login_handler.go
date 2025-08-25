package handler

import (
	"encoding/json"
	"net/http"
	"github.com/Flood-project/backend-flood/internal/login/usecase"
)

type LoginHandler struct {
	loginUseCase usecase.LoginManager
}

func NewLoginHandler(loginUseCase usecase.LoginManager) *LoginHandler{
	return &LoginHandler{
		loginUseCase: loginUseCase,
	}
}

func (handler *LoginHandler) Login(response http.ResponseWriter, request *http.Request) {
	var req struct {
		Email string `json:"email"`
		Password string `json:"password_hash"`
	}

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	account, err := handler.loginUseCase.Login(req.Email, req.Password)
	if err != nil {
		http.Error(response, "Email ou senha incorretos. Tente novamente.", http.StatusUnauthorized)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(response).Encode(&account)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}