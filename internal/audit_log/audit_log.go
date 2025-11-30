// internal/audit_log/audit_log.go
package auditlog

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Custom JSONB type para PostgreSQL
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
    if j == nil {
        return nil, nil
    }
    return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
    if value == nil {
        *j = nil
        return nil
    }
    
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("cannot scan %T into JSONB", value)
    }
    
    if len(bytes) == 0 {
        *j = nil
        return nil
    }
    
    var result map[string]interface{}
    if err := json.Unmarshal(bytes, &result); err != nil {
        return err
    }
    
    *j = result
    return nil
}

type AuditLog struct {
    ID           int64      `json:"id" db:"id"`
    TableName    string     `json:"table_name" db:"table_name"`
    RecordID     string     `json:"record_id" db:"record_id"`
    Operation    string     `json:"operation" db:"operation"`
    UserID       int32      `json:"user_id" db:"user_id"`
    UserEmail    string     `json:"user_email" db:"user_email"`
    OldData      JSONB      `json:"old_data" db:"old_data"`      // ✅ JSONB para PostgreSQL
    NewData      JSONB      `json:"new_data" db:"new_data"`      // ✅ JSONB para PostgreSQL
    ChangedFields []string  `json:"changed_fields" db:"changed_fields"`
    IPAddress    string     `json:"ip_address" db:"ip_address"`
    UserAgent    string     `json:"user_agent" db:"user_agent"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}