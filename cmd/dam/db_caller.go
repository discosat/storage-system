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
	fmt.Println("Logging rows in db_caller.go: ", rows)
	// Placeholder for Postgres service
	fmt.Println("Requesting image metadata from PostgrSQL using optimized query: ")
}

func MinioService() {
	// Placeholder for Minio service
	fmt.Println("Requesting images from Minio using optimized query: ")
}
