package usecase

import (
	"errors"

	accountRepository "github.com/Flood-project/backend-flood/internal/account_user/repository"
	"github.com/Flood-project/backend-flood/internal/token"
	tokenRepository "github.com/Flood-project/backend-flood/internal/token/repository"
	"golang.org/x/crypto/bcrypt"
)

type LoginManager interface {
	Login(email, password string) (string, string, error)
}

type loginUseCase struct {
	accountRepository accountRepository.AccountRepository
	token             token.TokenManager
	tokenRepository   tokenRepository.TokenRepository
}

func NewLogin(accountRepository accountRepository.AccountRepository, token token.TokenManager, tokenRepository tokenRepository.TokenRepository) LoginManager {
	return &loginUseCase{
		accountRepository: accountRepository,
		token:             token,
		tokenRepository:   tokenRepository,
	}
}

func (loginUseCase *loginUseCase) Login(email, password string) (string, string, error) {
	account, err := loginUseCase.accountRepository.GetByEmail(email)
	if err != nil {
		return "", "", errors.New("nenhum email encontrado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password_hash), []byte(password))
	if err != nil {
		return "", "", errors.New("senha inválida")
	}

	tokenString, refreshToken, err := loginUseCase.token.GenerateToken(*account)
	if err != nil {
		return "", "", errors.New("erro ao gerar token")
	}

	// tokenWithGroupUser := token.Token{
	// 	RowToken: tokenString,
	// 	Created: time.Now().Local(),
	// 	Expiration: time.Now().Add(time.Hour * 24).Local(),
	// 	IdAccount: account.Id_account,
	// }

	// err = loginUseCase.tokenRepository.Create(&tokenWithGroupUser)
	// if err != nil {
	// 	return "", err
	// }

	return tokenString, refreshToken, nil
}
