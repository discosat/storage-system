package observation

import (
	"github.com/discosat/storage-system/internal/Commands"
	"github.com/discosat/storage-system/internal/objectStore"
	"github.com/jmoiron/sqlx"
	"log"
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

func (p PsqlObservationRepository) CreateObservation(observationCommand Commands.ObservationCommand, metadata *ObservationMetadata) (int, error) {

	//tx := p.db.BeginTx()

	// TODO På et rollback ad SQL tx, skal billedet slettes
	objectReference, err := p.objectStore.SaveImage(observationCommand)
	if err != nil {
		return -1, err
	}

	var observationId int
	// TODO UserId
	err = p.db.QueryRow("INSERT INTO observation(observation_request_id, object_reference, user_id, bucket_name) VALUES ($1, $2, $3, $4) RETURNING id", observationCommand.ObservationRequestId, objectReference, observationCommand.UserId, observationCommand.Bucket).
		Scan(&observationId)
	if err != nil {
		log.Fatalf("err %v", err)
	}

	meta, err := p.CreateObservationMetadata(observationId, metadata)
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

func (p PsqlObservationRepository) CreateObservationMetadata(observationId int, metadata *ObservationMetadata) (int, error) {
	var metaId int
	err := p.db.QueryRow("INSERT INTO observation_metadata(observation_id, size, height, width, channels, timestamp, bits_pixels, image_offset, camera, location, gnss_date, gnss_time, gnss_speed, gnss_altitude, gnss_course) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,ST_SetSRID(ST_MakePoint($10,$11), 4326),$12,$13,$14,$15,$16) RETURNING id",
		observationId, metadata.Size, metadata.Height, metadata.Width, metadata.Channels, metadata.Timestamp, metadata.BitsPixels, metadata.ImageOffset, metadata.Camera, metadata.GnssLongitude, metadata.GnssLatitude, metadata.GnssDate, metadata.GnssTime, metadata.GnssSpeed, metadata.GnssAltitude, metadata.GnssCourse).
		Scan(&metaId)
	return metaId, err
}
