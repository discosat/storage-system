package dim

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStore struct {
	minioClient *minio.Client
}

func (m MinioStore) SaveFile(fileInfo *zip.File, openFile io.ReadCloser, bucketName string) (string, error) {
	status, err := m.minioClient.PutObject(context.Background(), bucketName, filepath.ToSlash(fileInfo.Name), openFile, fileInfo.FileInfo().Size(), minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("error in upload to minio: %v", err)
	}
	err = openFile.Close()
	if err != nil {
		return "", fmt.Errorf("error in closing file: %v", err)
	}
	return status.Key, nil
}

func (m MinioStore) SaveImage(fileInfo *multipart.FileHeader, openFile io.ReadCloser, bucketName string, observationName string) (string, error) {
	status, err := m.minioClient.PutObject(context.Background(), bucketName, filepath.ToSlash(observationName+"/"+fileInfo.Filename), openFile, fileInfo.Size, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("error in upload to minio: %v", err)
	}
	err = openFile.Close()
	if err != nil {
		return "", fmt.Errorf("error in closing file: %v", err)
	}
	return status.Key, nil
}

func (m MinioStore) BucketExists(bucketName string) (bool, error) {
	exists, err := m.minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return false, fmt.Errorf("unknown error ocurred %v", err)
	}
	return exists, nil
}

func NewMinioStore() MinioStore {
	er := godotenv.Load("cmd/dim/.env")
	if er != nil {
		log.Fatalf("NewMinioStore: Cant find env - %v", er)
	}
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	useSSL := true

	var err error
	minioC, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalf("NewMinioStore: %v", err)
	}

	return MinioStore{minioClient: minioC}
}
