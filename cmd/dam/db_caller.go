package dam

import (
	"database/sql"
	"fmt"
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

func PostgresService(query string, args []interface{}) {
	fmt.Println("Logging full Query in db_caller.go: ", query)
	fmt.Println("Logging Arguments in db_caller.go: ", args)
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal("An ERROR occured trying to create a query to the PostgreSQL database: ", err)
	}

	defer rows.Close()

	fmt.Println("Logging if we reach this")

	for rows.Next() {
		fmt.Println("Logging if we reach this2")
		var (
			id            int
			created_at    string
			updated_at    string
			measurementID int
			size          int
			height        int
			width         int
			channels      int
			timestamp     int64
			bitsPixels    int
			imageOffset   int
			camera        string
			location      string // We'll get to this in a moment
			gnssDate      int64
			gnssTime      int64
			gnssSpeed     float64
			gnssAltitude  float64
			gnssCourse    float64
		)

		fmt.Println("Logging to test if we get any data out?: ", size)

		err := rows.Scan(
			&id, &created_at, &updated_at, &measurementID, &size, &height, &width, &channels,
			&timestamp, &bitsPixels, &imageOffset, &camera, &location,
			&gnssDate, &gnssTime, &gnssSpeed, &gnssAltitude, &gnssCourse,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		fmt.Printf(
			"Row: ID=%d, Cam=%s, Loc=%s, Time=%d, Speed=%.2f, Alt=%.2f, Course=%.2f\n",
			id, camera, location, timestamp, gnssSpeed, gnssAltitude, gnssCourse,
		)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
	}
}

func MinioService() {
	// Placeholder for Minio service
	fmt.Println("Requesting images from Minio using optimized query: ")
}
