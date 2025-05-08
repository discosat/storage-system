package main

import (
	"github.com/discosat/storage-system/cmd/dim"
	"database/sql"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/observationRequest"
	"fmt"
	"github.com/discosat/storage-system/cmd/dam"
	"github.com/joho/godotenv"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
	"os"
)

func main() {
	fmt.SetFlags(log.LstdFlags | log.Lshortfile)
	slog.Info("starting DIM-DAM system backend")

	store := objectStore.NewMinioStore()
	db := ConfigDatabase()
	defer db.Close()

	//Initialize environment variables
	err := godotenv.Load("cmd/dam/.env")
	if err != nil {
		log.Fatalf("NewMinioStore: Cant find env - %v", err)
	}
	// Initialize the database connection
	dam.InitDB()

	// Initialize DIM and DAM services
	go dam.Start()
	go dim.Start(
		dim.NewDimController(
			dim.NewDimService(
				observation.NewPsqlObservationRepository(db, store),
				observationRequest.NewPsqlObservationRequestRepository(db),
			),
		),
	)

	slog.Info("DIM-DAM up and running")

	select {}
}

func ConfigDatabase() *sqlx.DB {
	db, err := sqlx.Open("pgx", fmt.Sprint("postgres://", os.Getenv("PGUSER"), ":", os.Getenv("PGPASSWORD"), "@", os.Getenv("PGHOST"), "/", os.Getenv("PGDATABASE")))
	if err != nil {
		log.Fatalf("Unable to configrue database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Connection to database could not be established: %v", err)
	}
	return db
}
