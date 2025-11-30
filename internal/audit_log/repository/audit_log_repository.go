// internal/audit_log/repository/audit_log_repository.go
package repository

import (
	"context"
	"fmt"
	"time"

	auditlog "github.com/Flood-project/backend-flood/internal/audit_log"
	"github.com/jmoiron/sqlx"
	//"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AuditLogManagement interface {
	Create(ctx context.Context, log *auditlog.AuditLog) error
	Fetch() ([]auditlog.AuditLog, error)
}

type auditLogManagement struct {
	DB *sqlx.DB
}

func NewAuditLogManagement(db *sqlx.DB) AuditLogManagement {
	return &auditLogManagement{
		DB: db,
	}
}

func (management *auditLogManagement) Create(ctx context.Context, log *auditlog.AuditLog) error {
    query := `
        INSERT INTO audit_logs (
            table_name, record_id, operation, user_id, user_email, 
            old_data, new_data, changed_fields, ip_address, user_agent, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id
    `
    
    // Converter para tipos que o PostgreSQL entende
    var oldData interface{} = nil
    var newData interface{} = nil
    
    if log.OldData != nil {
        oldData = log.OldData
    }
    if log.NewData != nil {
        newData = log.NewData
    }
    
    changedFieldsArray := make([]string, 0)
    if log.ChangedFields != nil {
        changedFieldsArray = log.ChangedFields
    }

    err := management.DB.QueryRowContext(
        ctx,
        query,
        log.TableName,
        log.RecordID,
        log.Operation,
        log.UserID,
        log.UserEmail,
        oldData,  // ✅ Pode ser nil
        newData,  // ✅ Pode ser nil  
        pq.Array(changedFieldsArray),
        log.IPAddress,
        log.UserAgent,
        time.Now(),
    ).Scan(&log.ID)
    
    if err != nil {
        fmt.Printf("❌ ERRO ao salvar audit log: %v\n", err)
        return err
    }
    
    return nil
}

func (management *auditLogManagement) Fetch() ([]auditlog.AuditLog, error) {
    query := `
        SELECT 
            id,
            COALESCE(table_name, '') as table_name,
            COALESCE(record_id, '') as record_id,
            COALESCE(operation, '') as operation, 
            COALESCE(user_id, 0) as user_id,
            COALESCE(user_email, '') as user_email,
            old_data,
            new_data, 
         
            COALESCE(ip_address, '') as ip_address,
            COALESCE(user_agent, '') as user_agent,
            created_at
        FROM audit_logs
        ORDER BY created_at DESC
    `

    var logs []auditlog.AuditLog
    err := management.DB.Select(&logs, query)
    if err != nil {
        fmt.Printf("Erro no Select: %v\n", err)
        return nil, err
    }

    fmt.Printf("Fetch retornou %d logs\n", len(logs))
    return logs, nil
}