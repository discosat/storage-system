package measurement

import (
	"database/sql"
	"github.com/discosat/storage-system/internal/objectStore"
	"mime/multipart"
)

type PsqlMeasurementRepository struct {
	db          *sql.DB
	objectStore objectStore.IDataStore
}

func NewPsqlMeasurementRepository(db *sql.DB, store objectStore.IDataStore) MeasurementRepository {
	return PsqlMeasurementRepository{db: db, objectStore: store}
}

func (p PsqlMeasurementRepository) GetById(id int) (Measurement, error) {
	//TODO implement me
	panic("implement me")
}

func (p PsqlMeasurementRepository) CreateMeasurement(file *multipart.FileHeader, bucket string, requestName string, observationId int, measurementRequestId int) (int, error) {

	oFile, err := file.Open()
	if err != nil {
		return -1, err
	}
	objectReference, err := p.objectStore.SaveImage(file, oFile, bucket, requestName)
	if err != nil {
		return -1, err
	}

	var measurementId int
	err = p.db.QueryRow("INSERT INTO measurement(object_reference, observation_id, measurement_request_id) VALUES ($1, $2, $3) RETURNING id", objectReference, observationId, measurementRequestId).
		Scan(&measurementId)
	return measurementId, err

}
