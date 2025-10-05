package repository

import (
	"github.com/Flood-project/backend-flood/internal/object_store"
	"github.com/jmoiron/sqlx"
)

type ObjectStoreManager interface {
	AddFile(file *object_store.FileData, fileByte []byte) error
	FetchFiles() ([]object_store.FileData, error)
}

type objectStoreManager struct {
	DB *sqlx.DB
}

func NewObjectStoreUseCase(db *sqlx.DB) ObjectStoreManager{
	return &objectStoreManager{
		DB: db,
	}
}

func (repository *objectStoreManager) AddFile(file *object_store.FileData, fileByte []byte) error {
	query := `INSERT INTO files (
		user_id, file_name, storage_key, url, size, content_type)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := repository.DB.QueryRow(
		query,
		file.UserID,
		file.FileName,
		file.StorageKey,
		file.URL,
		file.Size,
		file.ContentType,
	).Scan(&file.ID)
	if err != nil {
		return err 
	}

	return nil
}

func (repository *objectStoreManager) FetchFiles() ([]object_store.FileData, error) {
	query := `
		SELECT 
		  f.id,
		  f.user_id,
		  u.name,
		  u.email,
		  f.file_name,
		  f.storage_key,
		  f.url,
		  f.size,
		  f.content_type
		FROM files f
		INNER JOIN account u
		  ON u.id = f.user_id
	`
	var files []object_store.FileData

	err := repository.DB.Select(
		&files,
		query,
	)
	if err != nil {
		return nil, err
	}

	return files, nil
}