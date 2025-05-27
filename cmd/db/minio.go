package db

import (
	"bytes"
	"context"
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"os"
	"sync"
)

type MinIOClient struct {
	client *minio.Client
}

// Constructor
func NewMinIOClient() (*MinIOClient, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinIOClient{client: minioClient}, nil
}

func (mc *MinIOClient) FetchImages(imageInfo []interfaces.ImageMinIOData) ([]interfaces.RetrievedImages, error) {
	ctx := context.Background()
	results := make([]interfaces.RetrievedImages, 0)

	bucketGroups := make(map[string][]interfaces.ImageMinIOData)
	for _, image := range imageInfo {
		bucketGroups[image.BucketName] = append(bucketGroups[image.BucketName], image)
	}

	var wg sync.WaitGroup
	mu := sync.Mutex{}

	for bucket, images := range bucketGroups {
		wg.Add(1)
		go func(bucket string, images []interfaces.ImageMinIOData) {
			defer wg.Done()

			for _, image := range images {
				fmt.Println("Fetching from Bucket:", bucket, "Object:", image.ObjectReference)

				object, err := mc.client.GetObject(ctx, bucket, image.ObjectReference, minio.GetObjectOptions{})
				if err != nil {
					log.Printf("Failed to get object %s from bucket %s: %v", image.ObjectReference, image.BucketName, err)
					continue
				}

				buf := new(bytes.Buffer)
				if _, bufErr := io.Copy(buf, object); bufErr != nil {
					log.Printf("Failed to read object %s: %v", image.ObjectReference, bufErr)
					continue
				}

				mu.Lock()
				results = append(results, interfaces.RetrievedImages{
					ObjectReference: image.ObjectReference,
					BucketName:      image.BucketName,
					Image:           buf.Bytes(),
				})
				mu.Unlock()
			}
		}(bucket, images)
	}
	wg.Wait()

	return results, nil
}
