package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	//"strconv"
	"time"

	"github.com/Flood-project/backend-flood/internal/object_store"
	"github.com/Flood-project/backend-flood/internal/object_store/usecase"
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
	 //var fileEncoded object_store.FileData
	// err := json.NewDecoder(request.Body).Decode(&fileEncoded)
	// if err != nil {
	// 	log.Println("erro do decode: ", err)
	// 	http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }

	err := request.ParseMultipartForm(32 << 20) // 32MB max
    if err != nil {
        log.Println("Erro ao parse multipart form: ", err)
        http.Error(response, "Erro no formato do formulário", http.StatusBadRequest)
        return
    }

	userID := request.Context().Value("user_id").(int32)
	log.Println("id recebido: ", userID)

	// idWithCHi := chi.URLParam(request, "id")
	// chiID, err := strconv.Atoi(idWithCHi)
	// if err != nil {
	// 	log.Println("erro no id do chi: ", err, "chi id: ", chiID, "antes do strconv", idWithCHi)
	// 	log.Println("passing foor int32", int32(chiID))
	// 	http.Error(response, "Erro no id no usuário", http.StatusBadRequest)
	// 	return
	// }

	// userIDStr := request.FormValue("user_id")
	// userID, err := strconv.Atoi(userIDStr)
	// if err != nil {
	// 	log.Println("erro no id: ", err, "user id: ", userID, "antes do strconv", userIDStr)
	// 	log.Println("passing foor int32", int32(userID))
	// 	http.Error(response, "Erro no id no usuário", http.StatusBadRequest)
	// 	return
	// }

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
		UserID: userID,
		FileName: header.Filename,
		StorageKey: generateStorageKey(header.Filename),
		ContentType: header.Header.Get("Content-Type"),
		Size: header.Size,
		URL: url,
	}

	err = handler.usecase.AddFile(&fileData, fileBytes)
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