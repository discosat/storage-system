package main

import (
	"database/sql"
	"fmt"
	"github.com/discosat/storage-system/cmd/dim"
	"github.com/discosat/storage-system/internal/measurement"
	"github.com/discosat/storage-system/internal/measurementMetadata"
	"github.com/discosat/storage-system/internal/measurementRequest"
	"github.com/discosat/storage-system/internal/mission"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/request"
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
				request.NewPsqlRequestRepository(db),
				mission.NewPsqlMissionRepository(db),
				observation.NewPsqlObservationRepository(db),
				measurementRequest.NewPsqlMeasurementRequestRepository(db),
				measurement.NewPsqlMeasurementRepository(db, store),
				measurementMetadata.NewPsqlMeasurementMetadataRepository(db),
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
