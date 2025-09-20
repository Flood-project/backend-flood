package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Flood-project/backend-flood/internal/base"
	"github.com/Flood-project/backend-flood/internal/base/usecase"
	"github.com/go-chi/chi/v5"
)

type BaseHandler struct {
	baseUseCase usecase.BaseUseCase
}

func NewBaseHandler(baseUseCase usecase.BaseUseCase) *BaseHandler {
	return &BaseHandler{
		baseUseCase: baseUseCase,
	}
}

func (handler *BaseHandler) Create(response http.ResponseWriter, request *http.Request) {
	var base *base.Base

	err := json.NewDecoder(request.Body).Decode(&base)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.baseUseCase.Create(base)
	if err != nil {
		http.Error(response, "Erro ao adicionar novo tipo de base.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&base)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *BaseHandler) Fetch(response http.ResponseWriter, request *http.Request) {
	bases, err := handler.baseUseCase.Fetch()
	if err != nil {
		http.Error(response, "Erro ao listar tipos de base.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&bases)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *BaseHandler) Delete(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.baseUseCase.Delete(int32(id))
	if err != nil {
		http.Error(response, "Erro ao deletar tipo de base.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}