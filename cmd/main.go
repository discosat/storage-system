package main

import (
	"database/sql"
	"github.com/discosat/storage-system/cmd/dim"
	"github.com/discosat/storage-system/internal/flightPlan"
	"github.com/discosat/storage-system/internal/mission"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/observationMetadata"
	"github.com/discosat/storage-system/internal/observationRequest"
	"fmt"
	"github.com/discosat/storage-system/cmd/dam"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	fmt.Println("starting DIM-DAM system backend")

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
				flightPlan.NewPsqlFlightPlanRepository(db),
				mission.NewPsqlMissionRepository(db),
				observation.NewPsqlObservationRepository(db, store),
				observationRequest.NewPsqlObservationRequestRepository(db),
				observationMetadata.NewPsqlObservationMetadataRepository(db),
			),
		),
	)

	fmt.Println("DIM-DAM up and running")

	select {}
}

func ConfigDatabase() *sql.DB {
	db, err := sql.Open("pgx", fmt.Sprint("postgres://", os.Getenv("PGUSER"), ":", os.Getenv("PGPASSWORD"), "@", os.Getenv("PGHOST"), "/", os.Getenv("PGDATABASE")))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("ConfigDatabase: Cannot connect to database; %v", err)
	}
	return db
}
