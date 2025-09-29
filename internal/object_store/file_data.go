package object_store

type FileData struct {
	ID          int32  `json:"id" db:"id"`
	UserID      int32  `json:"user_id" db:"user_id"`
	UserName    string `json:"name" db:"name"`
	Email       string `json:"email" db:"email"`
	FileName    string `json:"file_name" db:"file_name"`
	StorageKey  string `json:"storage_key" db:"storage_key"`
	URL         string `json:"url" db:"url"`
	Size        int64  `json:"size" db:"size"`
	ContentType string `json:"content_type" db:"content_type"`
}