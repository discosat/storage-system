package observation

import (
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/jmoiron/sqlx"
	"log"
	"mime/multipart"
)

type PsqlObservationRepository struct {
	db          *sqlx.DB
	objectStore objectStore.IDataStore
}

func NewPsqlObservationRepository(db *sqlx.DB, store objectStore.IDataStore) ObservationRepository {
	return PsqlObservationRepository{db: db, objectStore: store}
}

func (p PsqlObservationRepository) GetObservation(id int) (Observation, error) {
	//TODO implement me
	panic("implement me")
}

func (p PsqlObservationRepository) CreateObservation(file *multipart.FileHeader, bucket string, flightPlanName string, observationRequest int) (int, error) {

	//tx := p.db.BeginTx()

	oFile, err := file.Open()
	if err != nil {
		return -1, err
	}
	// TODO På et rollback ad SQL tx, skal billedet slettes
	objectReference, err := p.objectStore.SaveImage(file, oFile, bucket, flightPlanName)
	if err != nil {
		return -1, err
	}

	var observationId int
	// TODO UserId
	err = p.db.QueryRow("INSERT INTO observation(observation_request_id, object_reference, user_id) VALUES ($1, $2, 1) RETURNING id", observationRequest, objectReference).
		Scan(&observationId)

	meta, err := p.CreateObservationMetadata(observationId, 10.4058633, 55.3821913)
	if err != nil {
		log.Fatalf("Går galt ved metadata upload: %v", err)
	}
	log.Println(meta)

	return observationId, err

}

func (p PsqlObservationRepository) GetObservationMetadata(id int) (ObservationMetadata, error) {
	//TODO implement me
	panic("implement me")
}

func (p PsqlObservationRepository) CreateObservationMetadata(observationId int, long float64, lat float64) (int, error) {
	var metaId int
	err := p.db.QueryRow("INSERT INTO observation_metadata(measurement_id, location) VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326)) RETURNING id", observationId, long, lat).
		Scan(&metaId)
	return metaId, err
}
