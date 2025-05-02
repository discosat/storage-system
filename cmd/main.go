package main

import (
	"database/sql"
	"fmt"
	"github.com/discosat/storage-system/cmd/dim"
	"github.com/discosat/storage-system/internal/flightPlan"
	"github.com/discosat/storage-system/internal/mission"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/observationMetadata"
	"github.com/discosat/storage-system/internal/observationRequest"
	"log"
	"os"
)

func main() {
	log.Println("starting DIM-DAM system backend")

	store := objectStore.NewMinioStore()
	db := ConfigDatabase()
	defer db.Close()

	//go dam.Start()
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

	log.Println("DIM-DAM up and running")

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
