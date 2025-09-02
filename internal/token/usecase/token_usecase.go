package usecase

import (
	"github.com/Flood-project/backend-flood/internal/token"
	"github.com/Flood-project/backend-flood/internal/token/repository"
)

type TokenUseCase interface {
	Create(token *token.Token) error
	Fetch() ([]token.Token, error)
}

type tokenUseCase struct {
	tokenRepository repository.TokenRepository
}

func (tokenUseCase *tokenUseCase) Create(token *token.Token) error {
	err := tokenUseCase.tokenRepository.Create(token)
	if err != nil {
		return err
	}

	return nil
}

func (tokenUseCase *tokenUseCase) Fetch() ([]token.Token, error) {
	token, err := tokenUseCase.tokenRepository.Fetch()
	if err != nil {
		return nil, err
	}

	return token, err
}