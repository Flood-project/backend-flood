package object_store

import "github.com/minio/minio-go/v7"

type MinIOConnectionResponse struct {
	Client *minio.Client
	Bucket string
}