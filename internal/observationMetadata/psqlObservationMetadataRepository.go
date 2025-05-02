package observationMetadata

import "database/sql"

type PsqlObservationMetadataRepository struct {
	db *sql.DB
}

func NewPsqlObservationMetadataRepository(db *sql.DB) ObservationMetadataRepository {
	return PsqlObservationMetadataRepository{db: db}
}

func (p PsqlObservationMetadataRepository) GetById(id int) (ObservationMetadata, error) {
	//TODO implement me
	panic("implement me")
}

func (p PsqlObservationMetadataRepository) CreateObservationMetadata(observationId int, long float64, lat float64) (int, error) {
	var metaId int
	err := p.db.QueryRow("INSERT INTO observation_metadata(measurement_id, location) VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326)) RETURNING id", observationId, long, lat).
		Scan(&metaId)
	return metaId, err
}
