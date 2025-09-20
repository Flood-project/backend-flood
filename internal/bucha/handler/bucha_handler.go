package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Flood-project/backend-flood/internal/bucha"
	"github.com/Flood-project/backend-flood/internal/bucha/usecase"
	"github.com/go-chi/chi/v5"
)

type BuchaHandler struct {
	buchaUseCase usecase.BuchaUseCase
}

func NewBuchaHandler(buchaUseCase usecase.BuchaUseCase) *BuchaHandler{
	return &BuchaHandler{
		buchaUseCase: buchaUseCase,
	}
}

func (handler *BuchaHandler) Create(response http.ResponseWriter, request *http.Request) {
	var bucha *bucha.Bucha

	err := json.NewDecoder(request.Body).Decode(&bucha)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.buchaUseCase.Create(bucha)
	if err != nil {
		http.Error(response, "Erro ao adicionar nova bucha.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&bucha)
	if err != nil {
		http.Error(response, "Erro ao enviar resposta para adicionar nova bucha.", http.StatusBadRequest)
		return 
	}
}

func (handler *BuchaHandler) Fetch(response http.ResponseWriter, request *http.Request) {
	buchas, err := handler.buchaUseCase.Fetch()
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&buchas)
	if err != nil {
		http.Error(response, "Erro ao enviar resposta para listar buchas.", http.StatusBadRequest)
		return 
	}
}

func (handler *BuchaHandler) Delete(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.buchaUseCase.Delete(int32(id))
	if err != nil {
		http.Error(response, "Erro ao deletar bucha.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}