package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Flood-project/backend-flood/internal/object_store"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type ObjectStoreManager interface {
	AddFile(file *object_store.FileData, fileByte []byte, productId int32) error
	FetchFiles() ([]object_store.FileData, error)
	GetFileUrl(storageKey string) (string, error)
	GetObject(storageKey string) ([]byte, string, error)
}

type objectStoreManager struct {
	DB *sqlx.DB
	Minio *object_store.MinIOConnectionResponse
}

func NewObjectStoreUseCase(db *sqlx.DB, minio *object_store.MinIOConnectionResponse) ObjectStoreManager {
	return &objectStoreManager{
		DB: db,
		Minio: minio,
	}
}

func (repository *objectStoreManager) AddFile(file *object_store.FileData, fileByte []byte, productId int32) error {
	query := `INSERT INTO files (
		product_id, file_name, storage_key, url, size, content_type)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := repository.DB.QueryRow(
		query,
		file.ProductID,
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
	reqParams := make(url.Values)
    reqParams.Set("response-content-type", "image/png")

	url, err := r.Minio.Client.PresignedGetObject(
		context.Background(),
		r.Minio.Bucket, 
		storageKey,
		time.Hour, 
		nil,
	)
	if err != nil {
		return "", err
	}

	urlStr := url.String()
	urlStr = strings.Replace(urlStr, "minio:9000", "localhost:9000", 1)
	log.Printf("URL gerada: %s", urlStr)
	return urlStr, nil
}

func (r *objectStoreManager) GetObject(storageKey string) ([]byte, string, error) {
	file, err := r.Minio.Client.GetObject(
		context.Background(),
		r.Minio.Bucket,
		storageKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, "", fmt.Errorf("erro ao buscar imagens no Minio: %w", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	stat, err := file.Stat()
	contentType := "image/png"
	if err != nil && stat.ContentType != "" {
		contentType = stat.ContentType
	}

	return fileBytes, contentType, nil
}
