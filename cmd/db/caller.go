package db

import (
	"github.com/discosat/storage-system/cmd/interfaces"
	"log"
)

type DBService struct {
	Postgres *PostgresClient
	MinIO    *MinIOClient
}

func NewDBService(pg *PostgresClient, minio *MinIOClient) *DBService {
	return &DBService{
		Postgres: pg,
		MinIO:    minio,
	}
}

func (svc *DBService) GetImages(query string, args []interface{}) ([]interfaces.RetrievedImages, error) {
	// 1. Query metadata from PostgreSQL
	metadata, err := svc.Postgres.PostgresQuery(query, args)
	if err != nil {
		log.Printf("Error querying PostgreSQL: %v", err)
		return nil, err
	}

	// 2. Prepare MinIO query payload
	var imageRequests []interfaces.ImageMinIOData
	for _, data := range metadata {
		imageRequests = append(imageRequests, interfaces.ImageMinIOData{
			BucketName:      data.BucketName,
			ObjectReference: data.ObjectReference,
		})
	}

	// 3. Query MinIO for image binaries
	images, err := svc.MinIO.MinioQuery(imageRequests)
	if err != nil {
		log.Printf("Error retrieving images from MinIO: %v", err)
		return nil, err
	}

	return images, nil
}
