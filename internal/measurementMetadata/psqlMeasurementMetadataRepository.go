package measurementMetadata

import "database/sql"

type PsqlMeasurementMetadataRepository struct {
	db *sql.DB
}

func NewPsqlMeasurementMetadataRepository(db *sql.DB) MeasurementMetadataRepository {
	return PsqlMeasurementMetadataRepository{db: db}
}

func (p PsqlMeasurementMetadataRepository) GetById(id int) (MeasurementMetadata, error) {
	//TODO implement me
	panic("implement me")
}

func (p PsqlMeasurementMetadataRepository) CreateMeasurementMetadata(measurementId int, long float64, lat float64) (int, error) {
	var metaId int
	err := p.db.QueryRow("INSERT INTO measurement_metadata(measurement_id, location) VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326)) RETURNING id", measurementId, long, lat).
		Scan(&metaId)
	return metaId, err
}
