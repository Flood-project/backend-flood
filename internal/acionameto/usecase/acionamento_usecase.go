package usecase

import (
	acionametos "github.com/Flood-project/backend-flood/internal/acionameto"
	"github.com/Flood-project/backend-flood/internal/acionameto/repository"
)

type AcionamentoUseCase interface {
	Create(acionamento *acionametos.Acionamento) error
	Fetch() ([]acionametos.Acionamento, error)
	Delete(id int32) error
	UpdateAcionamento(id int32, acionamento *acionametos.Acionamento) error
}

type acionamentoUseCase struct {
	acionamentoRepository repository.AcionamentoManagement
}

func NewAcionamentoUseCase(acionamentoRepository repository.AcionamentoManagement) AcionamentoUseCase {
	return &acionamentoUseCase{
		acionamentoRepository: acionamentoRepository,
	}
}

func (usecase *acionamentoUseCase) Create(acionamento *acionametos.Acionamento) error {
	err := usecase.acionamentoRepository.Create(acionamento)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *acionamentoUseCase) Fetch() ([]acionametos.Acionamento, error) {
	acionamentos, err := usecase.acionamentoRepository.Fetch()
	if err != nil {
		return nil, err
	}

	return acionamentos, nil
}

func (usecase *acionamentoUseCase) Delete(id int32) error {
	err := usecase.acionamentoRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *acionamentoUseCase) UpdateAcionamento(id int32, acionamento *acionametos.Acionamento) error {
	err := usecase.acionamentoRepository.UpdateAcionamento(id, acionamento)
	if err != nil {
		return err
	}

	return nil
}