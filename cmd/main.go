package main

import (
	"database/sql"
	"fmt"
	"github.com/discosat/storage-system/cmd/dam"
	"github.com/discosat/storage-system/cmd/dim"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/discosat/storage-system/internal/observation"
	"github.com/discosat/storage-system/internal/observationRequest"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	slog.Info("starting DIM-DAM system backend")

	err := godotenv.Load("cmd/dim/.env")
	if err != nil {
		log.Fatalf("NewMinioStore: Cant find env - %v", err)
	}

	store := objectStore.NewMinioStore()
	db := ConfigDatabase()
	defer db.Close()

	//Initialize environment variables
	//err = godotenv.Load("cmd/dam/.env")
	//if err != nil {
	//	log.Fatalf("NewMinioStore: Cant find env - %v", err)
	//}
	// Initialize the database connection
	//dam.InitDB()

	// Initialize DIM and DAM services
	damRouter, err := dam.ConfigureRouter()
	if err != nil {
		log.Fatal("Failed to start server")
	}
	go damRouter.Run(":8081")
	dimRouter := ConfigDimRouter(db, store)
	go dimRouter.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	slog.Info("DIM-DAM up and running")

	select {}
}

func ConfigDimRouter(db *sql.DB, store objectStore.IDataStore) *gin.Engine {
	return dim.ConfigureRouter(
		dim.NewDimController(
			dim.NewDimService(
				observation.NewPsqlObservationRepository(db, store),
				observationRequest.NewPsqlObservationRequestRepository(db),
			),
		),
	)
}

func ConfigDatabase() *sql.DB {
	db, err := sql.Open("pgx", fmt.Sprint("postgres://", os.Getenv("PGUSER"), ":", os.Getenv("PGPASSWORD"), "@", os.Getenv("PGHOST"), "/", os.Getenv("PGDATABASE")))
	if err != nil {
		log.Fatalf("Unable to configrue database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Connection to database could not be established: %v", err)
	}
	return db
}
