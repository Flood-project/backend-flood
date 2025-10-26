package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Flood-project/backend-flood/internal/product"
	"github.com/Flood-project/backend-flood/internal/product/usecase"
	"github.com/booscaaa/go-paginate/v3/paginate"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	productUseCase usecase.ProductUseCase
}

func NewProductHandler(productUseCase usecase.ProductUseCase) *ProductHandler {
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

func (handler *ProductHandler) FetchWithComponents(response http.ResponseWriter, request *http.Request) {
	productsWithComponents, err := handler.productUseCase.FetchWithComponents()
	if err != nil {
		log.Println(err, productsWithComponents)
		http.Error(response, "Erro ao listar produtos com componentes", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	err = json.NewEncoder(response).Encode(&productsWithComponents)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *ProductHandler) WithParams(w http.ResponseWriter, r *http.Request) {
	params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
	if err != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	rows, pageData, err := handler.productUseCase.WithParams(r.Context(), params)
	if err != nil {
		http.Error(w, "Não foi possível utlizar os filtros de pesquisa", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"products_with_params": rows,
		"total":                pageData.Total,
		"page":                 pageData.Page,
		"limit":                pageData.Limit,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *ProductHandler) GetByID(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	product, err := handler.productUseCase.GetByID(int32(id))
	if err != nil {
		http.Error(response, "Erro ao buscar produto por id", http.StatusInternalServerError)
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

func (handler *ProductHandler) Update(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var product product.Produt
	err = json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.productUseCase.Update(int32(id), &product)
	if err != nil {
		http.Error(response, "Erro ao atualizar produto", http.StatusInternalServerError)
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

func (handler *ProductHandler) Delete(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.productUseCase.Delete(int32(id))
	if err != nil {
		http.Error(response, "Erro ao deletar produto.", http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}
