package main

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const bucketRegion = "eu-north-0"

type MinioStore struct {
	minioClient *minio.Client
}

func (s MinioStore) SaveBatch(zipArchive *zip.ReadCloser, bucketName string) (string, error) {
	exists, err := s.minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return "", fmt.Errorf("SaveBatch: Unknown error ocurred %v", err)
	}

	if !exists {
		return "", fmt.Errorf("SaveBatch: No bucket with name '%v' exists. Specify an already existing bucket", bucketName)
	}

	for _, iFile := range zipArchive.File {
		if iFile.FileInfo().IsDir() {
			continue
		}
		oFile, _ := iFile.Open()
		log.Printf("SaveBatch: File name: %v", filepath.Clean(iFile.Name))

		status, err := s.minioClient.PutObject(context.Background(), bucketName, filepath.Clean(iFile.Name), oFile, iFile.FileInfo().Size(), minio.PutObjectOptions{})
		if err != nil {
			log.Printf("SaveBatch: Cannot upload thie file %v, error is %v", filepath.Base(iFile.Name), err)
			break
		}
		log.Println(status.Key)

	}
	return "success", nil
}

func newMinioStore() MinioStore {
	er := godotenv.Load()
	if er != nil {
		log.Fatalf("newMinioStore: Cant find env - %v", er)
	}
	//log.Println(os.Getenv("MINIO_ENDPONIT"))
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
		log.Fatalf("newMinioStore: %v", err)
	}

	return MinioStore{minioClient: minioC}
}
