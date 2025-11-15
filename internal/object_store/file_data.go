package object_store

type FileData struct {
	ID int32 `json:"id" db:"id"`
	ProductID int32  `json:"product_id" db:"product_id"`
	Codigo string `json:"codigo" db:"codigo"`
	//Email       string `json:"email" db:"email"`
	FileName    string `json:"file_name" db:"file_name"`
	StorageKey  string `json:"storage_key" db:"storage_key"`
	URL         string `json:"url" db:"url"`
	Size        int64  `json:"size" db:"size"`
	ContentType string `json:"content_type" db:"content_type"`
}
