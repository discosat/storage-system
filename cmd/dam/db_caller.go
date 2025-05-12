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
		log.Fatal("An ERROR occured trying to create a query to the PostgreSQL database: ", err)
	}

	defer rows.Close()

	var results []interfaces.ImageMetadata

	for rows.Next() {
		var row interfaces.ImageMetadata

		err := rows.Scan(
			&row.ID, &row.CreatedAt, &row.UpdatedAt, &row.MeasurementID, &row.Size, &row.Height, &row.Width, &row.Channels,
			&row.Timestamp, &row.BitsPixels, &row.ImageOffset, &row.Camera, &row.Location,
			&row.GnssDate, &row.GnssTime, &row.GnssSpeed, &row.GnssAltitude, &row.GnssCourse,
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

func MinioService() {
	// Placeholder for Minio service
	fmt.Println("Requesting images from Minio using optimized query: ")
}
