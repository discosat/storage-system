package db

import (
	"database/sql"
	"fmt"
	"github.com/discosat/storage-system/cmd/interfaces"
	"log"
)

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient(connStr string) (*PostgresClient, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("Postgres open error: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Postgres ping error: %w", err)
	}
	return &PostgresClient{db: db}, nil
}

func (pc *PostgresClient) FetchMetadata(query string, args []interface{}) ([]interfaces.ImageMetadata, error) {
	rows, err := pc.db.Query(query, args...)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("An ERROR occurred trying to create a query to the PostgreSQL database: %v", err)
		return nil, err
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
