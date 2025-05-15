package objectStore

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/discosat/storage-system/internal/Commands"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
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

func (m MinioStore) SaveImage(observationCommand Commands.ObservationCommand) (string, error) {
	status, err := m.minioClient.PutObject(context.Background(), observationCommand.Bucket, filepath.ToSlash(observationCommand.FlightPlanName+"/"+observationCommand.FileName), observationCommand.File, observationCommand.FileSize, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("error in upload to minio: %v", err)
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

func (m MinioStore) DeleteImage(imgRef string, bucketName string) (bool, error) {
	err := m.minioClient.RemoveObject(context.Background(), bucketName, imgRef, minio.RemoveObjectOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewMinioStore() MinioStore {
	err := godotenv.Load("cmd/dim/.env")
	if err != nil {
		log.Fatalf("NewMinioStore: Cant find env - %v", err)
	}
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
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
