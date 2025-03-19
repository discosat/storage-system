package dim

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
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
	err := godotenv.Load("cmd/dim/.env")
	if err != nil {
		log.Fatalf("NewMinioStore: Cant find env - %v", err)
	}
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
	useSSL := true

	minioC, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalf("NewMinioStore: %v", err)
	}

	// Check if minio is up and running
	_, err = minioC.ListBuckets(context.Background())
	if err != nil {
		log.Fatalf("Could not connect to Minio instance. Double check that it is up and running, and that you have provided correct credentials")
	}

	return MinioStore{minioClient: minioC}
}
