package config

import (
	"log"
	"os"
	"github.com/Flood-project/backend-flood/internal/object_store"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIO() (*object_store.MinIOConnectionResponse, error) {
	connect := object_store.MinIOConfig{
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		UseSSL:    false,
	}

	client, err := minio.New(connect.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(connect.AccessKey, connect.SecretKey, ""),
		Secure: connect.UseSSL,
	})

	if err != nil {
		log.Println("Erro na conexão com minIO")
		return nil, err
	}

	return &object_store.MinIOConnectionResponse{
		Client: client,
		Bucket: "files",
	}, nil
}
