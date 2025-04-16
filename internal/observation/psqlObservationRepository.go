package observation

import "database/sql"

type PsqlObservationRepository struct {
	db *sql.DB
}

func NewPsqlObservationRepository(db *sql.DB) ObservationRepository {
	return &PsqlObservationRepository{db: db}
}

func (p PsqlObservationRepository) GetById(id int) (Observation, error) {
	return Observation{}, nil
}

func (p PsqlObservationRepository) GetByRequest(requestId int) (Observation, error) {
	var observation Observation
	row := p.db.QueryRow("SELECT * FROM observation WHERE request_id = $1", requestId)
	err := row.Scan(&observation.Id, &observation.RequestId, &observation.UserId)
	return observation, err
}

func (p PsqlObservationRepository) CreateObservation(requestId int, userId int) (Observation, error) {
	var observation Observation
	err := p.db.QueryRow("INSERT INTO observation(request_id, user_id) VALUES ($1, $2) RETURNING id, request_id, user_id", requestId, userId).
		Scan(&observation.Id, &observation.RequestId, &observation.UserId)
	return observation, err
}
