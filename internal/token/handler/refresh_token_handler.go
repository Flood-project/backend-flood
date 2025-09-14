package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Flood-project/backend-flood/internal/account_user/usecase"
	"github.com/Flood-project/backend-flood/internal/token"
)

type TokenHandler struct {
	tokenManager token.TokenManager
	accountUseCase usecase.AccountUseCase
}

func NewTokenHandler(tokenManager token.TokenManager, accountUseCase usecase.AccountUseCase) TokenHandler{
	return TokenHandler{
		tokenManager: tokenManager,
		accountUseCase: accountUseCase,
	}
}

func (handler *TokenHandler) RefreshToken(response http.ResponseWriter, request *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil || req.RefreshToken == "" {
		http.Error(response, "Refresh token inválido", http.StatusBadRequest)
		return
	}

	claims, err := handler.tokenManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(response, "Refresh token inválido ou expirado.", http.StatusUnauthorized)
		return
	}

	account, err := handler.accountUseCase.GetByID(claims.IdUser)
	if err != nil {
		http.Error(response, "Conta ou usuário inválido.", http.StatusUnauthorized)
		return
	}	

	newToken, refreshToken, err := handler.tokenManager.GenerateToken(*account)
	if err != nil {
		http.Error(response, "Erro ao gerar novo token", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(response).Encode(map[string] string{
		"token": newToken,
		"refresh_token": refreshToken,
	})
	if err != nil {
		http.Error(response, "Erro ao enviar resposta do refresh token.", http.StatusBadRequest)
		return
	}
}