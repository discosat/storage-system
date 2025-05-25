package observation

import "github.com/discosat/storage-system/internal/Commands"

type ObservationRepository interface {
	GetObservation(id int) (Observation, error)
	CreateObservation(command Commands.CreateObservationCommand, metadata *ObservationMetadata) (int, error)
	GetObservationMetadata(id int) (ObservationMetadata, error)
	CreateObservationMetadata(observationId int, metadata *ObservationMetadata) (int, error)
}
