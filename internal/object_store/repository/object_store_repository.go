package repository

import (
	"context"
	"time"

	"github.com/Flood-project/backend-flood/internal/object_store"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStoreManager interface {
	AddFile(file *object_store.FileData, fileByte []byte) error
	FetchFiles() ([]object_store.FileData, error)
	GetFileUrl(storageKey string) (string, error)
}

type objectStoreManager struct {
	DB *sqlx.DB
}

func NewObjectStoreUseCase(db *sqlx.DB) ObjectStoreManager {
	return &objectStoreManager{
		DB: db,
	}
}

func (repository *objectStoreManager) AddFile(file *object_store.FileData, fileByte []byte) error {
	query := `INSERT INTO files (
		file_name, storage_key, url, size, content_type)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := repository.DB.QueryRow(
		query,
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
		  f.product_id,
		  p.codigo,
		  f.file_name,
		  f.storage_key,
		  f.url,
		  f.size,
		  f.content_type
		FROM files f
		INNER JOIN products p
		  ON p.id = f.product_id
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

func (r *objectStoreManager) GetFileUrl(storageKey string) (string, error) {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		return "", err
	}

	// Gera URL assinada que expira em 1 hora
	url, err := minioClient.PresignedGetObject(context.Background(), "files", storageKey, time.Hour, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
