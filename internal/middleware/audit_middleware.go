package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	auditlog "github.com/Flood-project/backend-flood/internal/audit_log"
	"github.com/Flood-project/backend-flood/internal/audit_log/usecase"
	"github.com/Flood-project/backend-flood/internal/token"
	"github.com/Flood-project/backend-flood/internal/token/util"
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
    case "GET":
        return "SELECT"
    case "POST":
        return "INSERT"
    case "PUT", "PATCH":
        return "UPDATE"
    case "DELETE":
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

func getClientIP(request *http.Request) string {
    if ip := request.Header.Get("X-Forwarded-For"); ip != "" {
        return strings.Split(ip, ",")[0]
    }
    if ip := request.Header.Get("X-Real-IP"); ip != "" {
        return ip
    }
    return request.RemoteAddr
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func (m *AuditMiddleware) GlobalAuditLog(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        wrappedWriter := &responseWriter{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }
        
        // Executa E DEPOIS salva automaticamente
        next.ServeHTTP(wrappedWriter, r)
        m.saveAuditLogAutomatically(r, wrappedWriter.statusCode, start) // CHAMADA AUTOMÁTICA
    })
}

func (m *AuditMiddleware) extractUserID(request *http.Request, statusCode int) int32 {
    // 1. Tenta pegar do context (se veio do middleware de autenticação)
    if ctxUserID := request.Context().Value("user_id"); ctxUserID != nil {
        if id, ok := ctxUserID.(int32); ok {
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
    fmt.Printf("=== GOROUTINE - Iniciando save no banco ===\n")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    fmt.Printf("Chamando usecase.Create...\n")
    if err := m.auditUseCase.Create(ctx, log); err != nil {
        fmt.Printf("❌ Erro ao salvar audit log automático: %v\n", err)
    }
}

func (m *AuditMiddleware) saveAuditLogAutomatically(r *http.Request, statusCode int, start time.Time) {
    duration := time.Since(start)
    
    fmt.Printf("=== MIDDLEWARE - Iniciando save automático ===\n")
    fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
    
    userID := m.extractUserID(r, statusCode)
    operation := getOperationFromMethod(r.Method)
    tableName := getTableFromPath(r.URL.Path)
    
    fmt.Printf("Dados extraídos: UserID=%d, Operation=%s, Table=%s\n", 
        userID, operation, tableName)
    
    auditLog := &auditlog.AuditLog{
        TableName: tableName,
        Operation: operation,
        UserID:    userID,
        IPAddress: getClientIP(r),
        UserAgent: r.UserAgent(),
        CreatedAt: time.Now(),
    }
    
    fmt.Printf("Criado auditLog: %+v\n", auditLog)
    
    go m.saveLogToDatabase(auditLog)
    
    log.Printf("AUDIT AUTO | UserID: %d | %s %s | Status: %d | Duration: %v",
        userID, r.Method, r.URL.Path, statusCode, duration)
        
    fmt.Printf("=== MIDDLEWARE - Save automático finalizado ===\n")
}
