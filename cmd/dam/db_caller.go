package dam

import (
	"database/sql"
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
	_ "github.com/lib/pq"
	"log"
	"os"
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

func MinIOService(imageInfo []interfaces.ImageMinIOData) {

	for _, image := range imageInfo {
		fmt.Println("Image ID:", image.ObjectReference)
		fmt.Println("Bucket Name:", image.BucketName)
	}
	/*
		endpoint := os.Getenv("MINIO_ENDPOINT")
		accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
		secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
		bucketRegion := os.Getenv("MINIO_BUCKET_REGION")

		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: false,
		})
		if err != nil {
			log.Fatal("Failed to create MinIO client:", err)
		}

		ctx := context.Background()

		for _, imageID := range imageIDs {
			objectName := fmt.Sprintf("image_%d.jpg", imageID)
			object, err := minioClient.GetObject(ctx, bucketRegion, objectName, minio.GetObjectOptions{})
			if err != nil {
				log.Println("Failed to get object from MinIO:", err)
				continue
			}

			buf := new(bytes.Buffer)
			_, bufErr := buf.ReadFrom(object)
			if bufErr != nil {
				log.Printf("Failed to read object %s: %v", objectName, err)
				continue
			}

			fmt.Printf("Successfully retrieved image: ", objectName)
		}

		// Placeholder for MinIO service
		fmt.Println("Requesting images from MinIO using Image ID: ", imageIDs)*/
}
