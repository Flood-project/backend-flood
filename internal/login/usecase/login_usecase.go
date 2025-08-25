package usecase

import (
	"errors"
	"github.com/Flood-project/backend-flood/internal/account_user/repository"
	"github.com/Flood-project/backend-flood/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type LoginManager interface {
	Login(email, password string) (string, error)
}

type loginUseCase struct {
	accountRepository repository.AccountRepository
	token             token.TokenManager
}

func NewLogin(accountRepository repository.AccountRepository, token token.TokenManager) LoginManager{
	return &loginUseCase{
		accountRepository: accountRepository,
		token: token,
	}
}

func (loginUseCase *loginUseCase) Login(email, password string) (string, error) {
	account, err := loginUseCase.accountRepository.GetByEmail(email)
	if err != nil {
		return "", errors.New("nenhum email encontrado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password_hash), []byte(password))
	if err != nil {
		return "", errors.New("senha inválida")
	}

	tokenString, err := loginUseCase.token.GenerateToken(*account)
	if err != nil {
		return "", errors.New("erro ao gerar token")
	}

	return tokenString, nil
}
