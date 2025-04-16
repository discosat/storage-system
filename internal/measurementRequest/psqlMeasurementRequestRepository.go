package measurementRequest

import "database/sql"

type PsqlMeasurementRequestRepository struct {
	db *sql.DB
}

func NewPsqlMeasurementRequestRepository(db *sql.DB) MeasurementRequestRepository {
	return &PsqlMeasurementRequestRepository{db: db}
}

func (p PsqlMeasurementRequestRepository) GetById(id int) (MeasurementRequest, error) {
	var measurementRequest MeasurementRequest
	err := p.db.QueryRow("SELECT * FROM measurement_request WHERE id = $1", id).Scan(&measurementRequest.Id, &measurementRequest.RId, &measurementRequest.MType)
	return measurementRequest, err
}
