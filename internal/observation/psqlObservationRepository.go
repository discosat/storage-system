package observation

import (
	"database/sql"
	"github.com/discosat/storage-system/internal/objectStore"
	"mime/multipart"
)

type PsqlObservationRepository struct {
	db          *sql.DB
	objectStore objectStore.IDataStore
}

func NewPsqlObservationRepository(db *sql.DB, store objectStore.IDataStore) ObservationRepository {
	return PsqlObservationRepository{db: db, objectStore: store}
}

func (p PsqlObservationRepository) GetById(id int) (Observation, error) {
	//TODO implement me
	panic("implement me")
}

func (p PsqlObservationRepository) CreateObservation(file *multipart.FileHeader, bucket string, flightPlanName string, observationRequest int) (int, error) {

	oFile, err := file.Open()
	if err != nil {
		return -1, err
	}
	objectReference, err := p.objectStore.SaveImage(file, oFile, bucket, flightPlanName)
	if err != nil {
		return -1, err
	}

	var measurementId int
	// TODO UserId
	err = p.db.QueryRow("INSERT INTO observation(observation_request_id, object_reference, user_id) VALUES ($1, $2, 1) RETURNING id", observationRequest, objectReference).
		Scan(&measurementId)
	return measurementId, err

}
