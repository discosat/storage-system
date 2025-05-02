package observationRequest

import "database/sql"

type PsqlObservationRequestRepository struct {
	db *sql.DB
}

func NewPsqlObservationRequestRepository(db *sql.DB) ObservationRequestRepository {
	return &PsqlObservationRequestRepository{db: db}
}

func (p PsqlObservationRequestRepository) GetById(id int) (ObservationRequest, error) {
	var observationRequest ObservationRequest
	err := p.db.QueryRow("SELECT * FROM observation_request WHERE id = $1", id).Scan(&observationRequest.Id, &observationRequest.CreatedAt, &observationRequest.UpdatedAt, &observationRequest.RId, &observationRequest.OType)
	return observationRequest, err
}
