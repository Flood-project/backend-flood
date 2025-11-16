package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	acionametos "github.com/Flood-project/backend-flood/internal/acionameto"
	"github.com/Flood-project/backend-flood/internal/acionameto/usecase"
	"github.com/go-chi/chi/v5"
)

type AcionamentoHandler struct {
	acionamentoUseCase usecase.AcionamentoUseCase
}

func NewAcionamentoHandler(acionamentoUseCase usecase.AcionamentoUseCase) *AcionamentoHandler{
	return &AcionamentoHandler{
		acionamentoUseCase: acionamentoUseCase,
	}
}

func (handler *AcionamentoHandler) Create(response http.ResponseWriter, request *http.Request) {
	var acionamento *acionametos.Acionamento

	err := json.NewDecoder(request.Body).Decode(&acionamento)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.acionamentoUseCase.Create(acionamento)
	if err != nil {
		http.Error(response, "Erro ao adicionar novo tipo de acionamento.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&acionamento)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *AcionamentoHandler) Fetch(response http.ResponseWriter, request *http.Request) {
	acionamentos, err := handler.acionamentoUseCase.Fetch()
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&acionamentos)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *AcionamentoHandler) Delete(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.acionamentoUseCase.Delete(int32(id))
	if err != nil {
		http.Error(response, "Erro ao deletar acionamento.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}

func (handler *AcionamentoHandler) UpdateAcionamento(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, "Acionamento não encontrado. ", http.StatusBadRequest)
		return
	}

	var acionamento acionametos.Acionamento

	err = json.NewDecoder(request.Body).Decode(&acionamento)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.acionamentoUseCase.UpdateAcionamento(int32(id), &acionamento)
	if err != nil {
		http.Error(response, "Erro ao editar acionamento.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&acionamento)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}