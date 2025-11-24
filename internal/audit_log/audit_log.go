package auditlog

import "time"

type AuditLog struct {
	ID            int       `json:"id" db:"id"`
	TableName     string    `json:"table_name" db:"table_name"`
	Operation     string    `json:"operation" db:"operation"` 
	UserID        int32     `json:"user_id" db:"user_id"`  
	UserEmail     string    `json:"user_email" db:"user_email"`
	OldData       []byte    `json:"old_data" db:"old_data"`       
	NewData       []byte    `json:"new_data" db:"new_data"`       
	ChangedFields []string  `json:"changed_fields" db:"changed_fields"`
	IPAddress     string    `json:"ip_address" db:"ip_address"`
	UserAgent     string    `json:"user_agent" db:"user_agent"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
