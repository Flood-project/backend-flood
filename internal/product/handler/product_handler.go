package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Flood-project/backend-flood/internal/product"
	"github.com/Flood-project/backend-flood/internal/product/usecase"
)

type ProductHandler struct {
	productUseCase usecase.ProductUseCase
}

func NewProductHandler(productUseCase usecase.ProductUseCase) *ProductHandler{
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

func (handler *ProductHandler) Create(response http.ResponseWriter, request *http.Request) {
	var product product.Produt
	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	err = handler.productUseCase.Create(&product)
	if err != nil {
		http.Error(response, "Erro ao adicionar novo produto.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&product)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return 
	}
}

func (handler *ProductHandler) Fetch(response http.ResponseWriter, request *http.Request) {
	products, err := handler.productUseCase.Fetch()
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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