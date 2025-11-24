package repository

import (
	"context"
	"fmt"
	"time"

	auditlog "github.com/Flood-project/backend-flood/internal/audit_log"
	"github.com/jmoiron/sqlx"
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
            table_name, operation, user_id, user_email, old_data, 
            new_data, changed_fields, ip_address, user_agent, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `
    
    changedFieldsArray := make([]string, 0)
    if log.ChangedFields != nil {
        changedFieldsArray = log.ChangedFields
    }

	var oldData interface{} = nil
    var newData interface{} = nil
    
    if len(log.OldData) > 0 {
        oldData = log.OldData
    }
    
    if len(log.NewData) > 0 {
        newData = log.NewData
    }
    
    err := management.DB.QueryRowContext(
        ctx,
        query,
        log.TableName,
        log.Operation,
        log.UserID,
        log.UserEmail,
        oldData,
        newData,
        pq.Array(changedFieldsArray), // usando github.com/lib/pq
        log.IPAddress,
        log.UserAgent,
        time.Now(),
    ).Scan(&log.ID)
	if err != nil {
        fmt.Printf("❌ ERRO ao salvar audit log: %v\n", err)
    }

    return err
}

func (management *auditLogManagement) Fetch() ([]auditlog.AuditLog, error) {
	 query := `
        SELECT 
            id,
            COALESCE(table_name, '') as table_name,
            COALESCE(operation, '') as operation, 
            COALESCE(user_id, 0) as user_id,
            COALESCE(user_email, '') as user_email,
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

