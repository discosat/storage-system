package main

import (
	"archive/zip"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"path/filepath"
)

const bucketName = "disco2data"
const bucketRegion = "eu-north-0"

type MinioStore struct {
	minioClient *minio.Client
}

func (s MinioStore) SaveBatch(zipArchive *zip.ReadCloser) (string, error) {
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
	endpoint := "app-disco-minio.cloud.sdu.dk"
	accessKeyID := ""
	secretAccessKey := ""
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
