package usecase

import (
	"bytes"
	"context"

	"log"

	"github.com/Flood-project/backend-flood/internal/object_store"
	"github.com/Flood-project/backend-flood/internal/object_store/repository"
	"github.com/minio/minio-go/v7"
)

type ObjectStoreUseCase interface {
	AddFile(file *object_store.FileData, fileByte []byte, productId int32) error
	FetchFiles() ([]object_store.FileData, error)
	GetFileUrl(storageKey string) (string, error)
	GetObject(storageKey string) ([]byte, string, error)
}

type objectStoreUseCase struct {
	objectStoreRepository   repository.ObjectStoreManager
	minIOConnectionResponse object_store.MinIOConnectionResponse
}

func NewObjectStoreUseCase(objectStoreRepository repository.ObjectStoreManager, minIOConnectionResponse object_store.MinIOConnectionResponse) ObjectStoreUseCase {
	return &objectStoreUseCase{
		objectStoreRepository:   objectStoreRepository,
		minIOConnectionResponse: minIOConnectionResponse,
	}
}

func (usecase *objectStoreUseCase) AddFile(file *object_store.FileData, fileByte []byte, productId int32) error {
	ctx := context.Background()

	exists, err := usecase.minIOConnectionResponse.Client.BucketExists(ctx, usecase.minIOConnectionResponse.Bucket)
	if err != nil {
		log.Println("Bucket não existe. ", err)
		return err
	}

	if !exists {
		err = usecase.minIOConnectionResponse.Client.MakeBucket(ctx, usecase.minIOConnectionResponse.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Println("Erro ao criar bucket. ", err)
			return err
		}
	}

	reader := bytes.NewReader(fileByte)

	_, err = usecase.minIOConnectionResponse.Client.PutObject(
		ctx,
		usecase.minIOConnectionResponse.Bucket,
		file.StorageKey,
		reader,
		int64(len(fileByte)),
		minio.PutObjectOptions{
			ContentType: file.ContentType,
		},
	)

	if err != nil {
		log.Println("Erro ao salvar imagem no minIO. ", err)
		return err
	}
	return usecase.objectStoreRepository.AddFile(file, fileByte, productId)
	// file.URL = fmt.Sprintf("/%s/%s", usecase.minIOConnectionResponse.Bucket, info.Key)
	// file.Size = info.Size

	// err = usecase.objectStoreRepository.AddFile(file, fileByte)
	// if err != nil {
	// 	log.Println("Erro ao salvar arquivo no banco de dados. ", err)
	// 	usecase.minIOConnectionResponse.Client.RemoveObjects(ctx, usecase.minIOConnectionResponse.Bucket, make(<-chan minio.ObjectInfo), minio.RemoveObjectsOptions{})
	// 	return err
	// }

	// log.Println("arquivo salvo: ", file.FileName, file.Size)
	// return nil
}

func (usecase *objectStoreUseCase) FetchFiles() ([]object_store.FileData, error) {
	files, err := usecase.objectStoreRepository.FetchFiles()
	if err != nil {
		log.Println("erro ao retornar arquivos: ", err)
		return nil, err
	}

	return files, err
}

func (uc *objectStoreUseCase) GetFileUrl(storageKey string) (string, error) {
	url, err := uc.objectStoreRepository.GetFileUrl(storageKey)
	if err != nil {
		log.Println("erro arquivos: ", err)
		return "", err
	}
	return url, nil
}


func (uc *objectStoreUseCase) GetObject(storageKey string) ([]byte, string, error){
	fileByte, url, err := uc.objectStoreRepository.GetObject(storageKey)
	if err != nil {
		log.Println("Erro ao buscar imagem", err)
		return nil, "", err
	}
	return fileByte, url, nil
}