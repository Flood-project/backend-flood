package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	auditlog "github.com/Flood-project/backend-flood/internal/audit_log"
	"github.com/Flood-project/backend-flood/internal/audit_log/usecase"
	"github.com/Flood-project/backend-flood/internal/token"
	"github.com/Flood-project/backend-flood/internal/token/util"
	"github.com/go-chi/chi/v5"
)

type AuditMiddleware struct {
    auditUseCase usecase.AuditLogUseCase
	tokenManager token.TokenManager
}

func NewAuditMiddleware(auditUseCase usecase.AuditLogUseCase, tokenManager token.TokenManager) *AuditMiddleware {
    return &AuditMiddleware{
        auditUseCase: auditUseCase,
		tokenManager: tokenManager,
    }
}

// Helper functions (mantenha as mesmas)
func getOperationFromMethod(method string) string {
    switch method {
    case http.MethodPost:
        return "INSERT"
    case http.MethodPut, http.MethodPatch:
        return "UPDATE"
    case http.MethodDelete:
        return "DELETE"
    default:
        return "UNKNOWN"
    }
}

func getTableFromPath(path string) string {
    paths := strings.Split(path, "/")
    if len(paths) >= 3 {
        return paths[2]
    }
    return "unknown"
}

// func getClientIP(request *http.Request) string {
//     if ip := request.Header.Get("X-Forwarded-For"); ip != "" {
//         return strings.Split(ip, ",")[0]
//     }
//     if ip := request.Header.Get("X-Real-IP"); ip != "" {
//         return ip
//     }
//     return request.RemoteAddr
// }

type responseWriter struct {
    http.ResponseWriter
    body *bytes.Buffer
    statusCode int
}

func (r *responseWriter) Write(b []byte) (int, error) {
    if r.body != nil {
        r.body.Write(b)
    }
	return r.ResponseWriter.Write(b)
}


func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func (m *AuditMiddleware) GlobalAuditLog(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !m.shouldAudit(r.Method) {
			next.ServeHTTP(w, r)
			return
		}

        start := time.Now()

        var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

        //r.body.Write(b)r.body.Write(b)r.body.Write(b)r.body.Write(b)
        wrappedWriter := &responseWriter{
            ResponseWriter: w,
            body: &bytes.Buffer{},
            statusCode:     http.StatusOK,
        }
        
        // Executa E DEPOIS salva automaticamente
        next.ServeHTTP(wrappedWriter, r)
        if wrappedWriter.statusCode >= 200 && wrappedWriter.statusCode < 300 {
			go m.saveAuditLogAutomatically(r, requestBody, wrappedWriter.body.Bytes(), wrappedWriter.statusCode, start)
		}
    })
}

func (m *AuditMiddleware) shouldAudit(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete
}

func (m *AuditMiddleware) extractTableInfo(r *http.Request) (string, string) {
	// Extrair nome da tabela da URL
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	
	if len(pathParts) == 0 {
		return "unknown", ""
	}

	tableName := pathParts[0]
	var recordID string

	// Tentar extrair ID do registro
	if len(pathParts) >= 2 {
		lastPart := pathParts[len(pathParts)-1]
		if !isNumeric(lastPart) {
			recordID = lastPart
		}
	}

	// Se não encontrou, tentar pelo chi URL params
	if recordID == "" {
		if id := chi.URLParam(r, "id"); id != "" {
			recordID = id
		}
	}

	return tableName, recordID
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (m *AuditMiddleware) extractDataChanges(r *http.Request, requestBody, responseBody []byte, operation string) (json.RawMessage, json.RawMessage, []string) {
	var oldData, newData json.RawMessage
	var changedFields []string

	switch operation {
	case "INSERT":
		if len(requestBody) > 0 {
			newData = requestBody
		}
	case "UPDATE":
		// Request body contém os novos dados
		if len(requestBody) > 0 {
			newData = requestBody
			
			// Extrair campos alterados do JSON
			var data map[string]interface{}
			if err := json.Unmarshal(requestBody, &data); err == nil {
				for field := range data {
					changedFields = append(changedFields, field)
				}
			}
		}
	case "DELETE":
		// Response body pode conter dados do registro deletado
		if len(responseBody) > 0 {
			oldData = responseBody
		}
	}

	return oldData, newData, changedFields
}

func (m *AuditMiddleware) extractUserID(request *http.Request, statusCode int) int32 {
    // 1. Tenta pegar do context (se veio do middleware de autenticação)
    if ctxUserID := request.Context().Value("user_id"); ctxUserID != nil {
        if id, ok := ctxUserID.(int32); ok {
            return id
        }
    }

    if ctxUserID := request.Context().Value("email"); ctxUserID != nil {
        if id, ok := ctxUserID.(int32); ok {
            log.Println(id)
            return id
        }
    }

    authHeader := request.Header.Get("Authorization")
    if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
        accessTokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        token, err := m.tokenManager.ValidateToken(accessTokenString)
        if err == nil && token.Type == "access" {
            claims := util.ExtractClaims(token)
            return claims.IdUser
        }
    }

    return 0
}

func (m *AuditMiddleware) saveLogToDatabase(log *auditlog.AuditLog) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    fmt.Printf("Chamando usecase.Create...\n")
    if err := m.auditUseCase.Create(ctx, log); err != nil {
        fmt.Printf("❌ Erro ao salvar audit log automático: %v\n", err)
    }
}

func getClientIP(r *http.Request) string {
	headers := []string{"X-Real-Ip", "X-Forwarded-For", "CF-Connecting-IP"}
	for _, header := range headers {
		if ip := r.Header.Get(header); ip != "" {
			return strings.Split(ip, ",")[0]
		}
	}
	return r.RemoteAddr
}

func (m *AuditMiddleware) saveAuditLogAutomatically(r *http.Request, requestBody, responseBody []byte, statusCode int, start time.Time) {
    ctx := r.Context()

    userID := m.extractUserID(r, statusCode)
    operation := getOperationFromMethod(r.Method)
    tableName, recordID := m.extractTableInfo(r)

    if ctx.Value("audit_processed") != nil {
        return
    }
    ctx = context.WithValue(r.Context(), "audit_processed", true)
    r = r.WithContext(ctx)
    
    
    // ✅ CAPTURAR DADOS REAIS
    _, newData, _ := m.extractRealData(r, requestBody, responseBody, operation, tableName)

    auditLog := &auditlog.AuditLog{
        TableName:     tableName,
        RecordID:      recordID,
        Operation:     operation,
        UserID:        userID,
   
        OldData:       nil,        // ✅ Dados ANTES
        NewData:       newData,        // ✅ Dados DEPOIS  
        ChangedFields: nil,  // ✅ Campos alterados
        IPAddress:     getClientIP(r),
        UserAgent:     r.UserAgent(),
        CreatedAt:     time.Now(),
    }
    
    fmt.Printf("✅ Salvando audit log COMPLETO: Table=%s, Record=%s, Operation=%s\n", 
        tableName, recordID, operation)
    
    go m.saveLogToDatabase(auditLog)
}

func (m *AuditMiddleware) extractRealData(r *http.Request, requestBody, responseBody []byte, operation, tableName string) (auditlog.JSONB, auditlog.JSONB, []string) {
    var oldData, newData auditlog.JSONB
    var changedFields []string

    switch operation {
    case "INSERT":
        // Para INSERT, o request body tem os novos dados
        if len(requestBody) > 0 {
            newData = m.parseRequestBody(requestBody)
        }
        
    case "UPDATE":
        // Para UPDATE, precisamos capturar:
        // - Dados antigos (do banco - isso é mais complexo)
        // - Dados novos (do request)
        // - Campos alterados
        
        if len(requestBody) > 0 {
            newData = m.parseRequestBody(requestBody)
            
            // Extrair campos alterados do JSON
            if newData != nil {
                for field := range newData {
                    changedFields = append(changedFields, field)
                }
            }
        }
        
        // 💡 DICA: Para capturar dados antigos, você precisaria:
        // 1. Fazer uma query antes da alteração para pegar o estado anterior
        // 2. Ou seu handler poderia fornecer esses dados via context
        
    case "DELETE":
        // Para DELETE, podemos tentar capturar dados do response
        if len(responseBody) > 0 {
            oldData = m.parseResponseBody(responseBody)
        }
    }

    return oldData, newData, changedFields
}

func (m *AuditMiddleware) parseRequestBody(body []byte) auditlog.JSONB {
    var data map[string]interface{}
    if err := json.Unmarshal(body, &data); err == nil {
        return auditlog.JSONB(data)
    }
    return nil
}

func (m *AuditMiddleware) parseResponseBody(body []byte) auditlog.JSONB {
    var data map[string]interface{}
    if err := json.Unmarshal(body, &data); err == nil {
        return auditlog.JSONB(data)
    }
    return nil
}
