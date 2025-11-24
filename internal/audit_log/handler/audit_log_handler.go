package handler

import (
	"encoding/json"
	"log"
	"net/http"

	auditlog "github.com/Flood-project/backend-flood/internal/audit_log"
	"github.com/Flood-project/backend-flood/internal/audit_log/usecase"
)

type AuditLogHandler struct {
	usecase usecase.AuditLogUseCase
}

func NewBaseHandler(usecase usecase.AuditLogUseCase) *AuditLogHandler {
	return &AuditLogHandler{
		usecase: usecase,
	}
}

func (handler *AuditLogHandler) Create(response http.ResponseWriter, request *http.Request) {
	var log auditlog.AuditLog
    
    if err := json.NewDecoder(request.Body).Decode(&log); err != nil {
        response.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(response).Encode(map[string]string{"error": "invalid request body"})
        return
    }
    
    if err := handler.usecase.Create(request.Context(), &log); err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(response).Encode(map[string]string{"error": "failed to create log"})
        return
    }
    
    response.WriteHeader(http.StatusCreated)
    json.NewEncoder(response).Encode(log)
}

func (handler *AuditLogHandler) Fetch(response http.ResponseWriter, request *http.Request) {
	logs, err := handler.usecase.Fetch()
	if err != nil {
		log.Println("error to fetch logs", err)
		http.Error(response, "Erro ao listar logs.", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	err = json.NewEncoder(response).Encode(&logs)
	if err != nil {
		http.Error(response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}