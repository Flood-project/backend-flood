package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	//"strconv"
	"time"

	"github.com/Flood-project/backend-flood/internal/object_store"
	"github.com/Flood-project/backend-flood/internal/object_store/usecase"
	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5"
)

type ObjectStoreHandler struct {
	usecase usecase.ObjectStoreUseCase
}

func NewObjectStoreHandler(usecase usecase.ObjectStoreUseCase) *ObjectStoreHandler {
	return &ObjectStoreHandler{
		usecase: usecase,
	}
}

func (handler *ObjectStoreHandler) Create(response http.ResponseWriter, request *http.Request) {
	// productIdStr := chi.URLParam(request, "id")
	// log.Println("id: ", productIdStr)
	// productId, err := strconv.Atoi(productIdStr)
	// if err != nil {
	// 	http.Error(response, "ID de produto não encontrado. ", http.StatusBadRequest)
	// 	return
	// } //no need productId

	//var fileEncoded object_store.FileData
	// err := json.NewDecoder(request.Body).Decode(&fileEncoded)
	// if err != nil {
	// 	log.Println("erro do decode: ", err)
	// 	http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }

	// err := request.ParseMultipartForm(32 << 20) // 32MB max
	// if err != nil {
	// 	log.Println("Erro ao parse multipart form: ", err)
	// 	http.Error(response, "Erro no formato do formulário", http.StatusBadRequest)
	// 	return
	// }

	//productID := request.Context().Value("product_id").(int32)
	//log.Println("id recebido: ", productID)

	// idWithCHi := chi.URLParam(request, "id")
	// chiID, err := strconv.Atoi(idWithCHi)
	// if err != nil {
	// 	log.Println("erro no id do chi: ", err, "chi id: ", chiID, "antes do strconv", idWithCHi)
	// 	log.Println("passing foor int32", int32(chiID))
	// 	http.Error(response, "Erro no id no usuário", http.StatusBadRequest)
	// 	return
	// }

	productIDStr := chi.URLParam(request, "product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(response, "Erro ao receber ID do produto", http.StatusBadRequest)
		return
	}

	file, header, err := request.FormFile("file")
	if err != nil {
		log.Println("formato: ", err)
		http.Error(response, "Erro no formato do arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(response, "Erro ao processar arquivo", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/%s", header.Filename)
	fileData := object_store.FileData{
		ProductID:   int32(productID),
		FileName:    header.Filename,
		StorageKey:  generateStorageKey(header.Filename),
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
		URL:         url,
	}

	err = handler.usecase.AddFile(&fileData, fileBytes, int32(productID))
	if err != nil {
		log.Println(err)
		http.Error(response, "Erro ao salvar arquivo no banco de dados", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(fileData)
	if err != nil {
		log.Println("erro do encode: ", err)
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func generateStorageKey(filename string) string {
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
}

func (handler *ObjectStoreHandler) Fetch(response http.ResponseWriter, request *http.Request) {
	files, err := handler.usecase.FetchFiles()
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&files)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (handler *ObjectStoreHandler) GetFileUrl(response http.ResponseWriter, request *http.Request) {
	storageKey := chi.URLParam(request, "storageKey")

	url, err := handler.usecase.GetFileUrl(storageKey)
	if err != nil {
		http.Error(response, "Erro ao gerar URL", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]string{"url": url})
}

func (handler *ObjectStoreHandler) ServeImage(response http.ResponseWriter, request *http.Request) {
	storageKey := chi.URLParam(request, "storageKey")
	if storageKey == "" {
		http.Error(response, "storageKey é obrigatório", http.StatusBadRequest)
		return
	}

	fileBytes, contentType, err := handler.usecase.GetObject(storageKey)
	if err != nil {
		log.Printf("Erro ao buscar imagem %s: %v", storageKey, err)
		http.Error(response, "Imagem não encontrada", http.StatusNotFound)
		return
	}

	response.Header().Set("Content-Type", contentType)
	response.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	response.Header().Set("Cache-Control", "public, max-age=3600")

	response.Write(fileBytes)
}
