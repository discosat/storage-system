package dam

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"os"
	"sync"
)

var db *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("PGHOST"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
}

func PostgresService(query string, args []interface{}) ([]interfaces.ImageMetadata, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatal("An ERROR occured trying to create a query to the PostgreSQL database: ", err)
	}

	defer rows.Close()

	var results []interfaces.ImageMetadata

	for rows.Next() {
		var row interfaces.ImageMetadata

		err := rows.Scan(
			&row.ID, &row.CreatedAt, &row.UpdatedAt, &row.ObservationID, &row.Size, &row.Height, &row.Width, &row.Channels,
			&row.Timestamp, &row.BitsPixels, &row.ImageOffset, &row.Camera, &row.Location,
			&row.GnssDate, &row.GnssTime, &row.GnssSpeed, &row.GnssAltitude, &row.GnssCourse, &row.BucketName, &row.ObjectReference,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
	}

	return results, nil
}

func MinIOService(imageInfo []interfaces.ImageMinIOData) ([]interfaces.RetrievedImages, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal("Failed to create MinIO client:", err)
	}

	/*_, err = minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalf("Could not connect to Minio instance. Double check that it is up and running, and that you have provided correct credentials")
	}*/

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

				object, err := minioClient.GetObject(ctx, image.BucketName, image.ObjectReference, minio.GetObjectOptions{})
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

	return results, err
}
