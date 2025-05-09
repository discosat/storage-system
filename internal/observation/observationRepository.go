package observation

import "github.com/discosat/storage-system/internal/Commands"

type ObservationRepository interface {
	GetObservation(id int) (Observation, error)
	// CreateObservation TODO Should preferably not depend on specific file type. Check if a Reader could do
	CreateObservation(command Commands.ObservationCommand) (int, error)

	GetObservationMetadata(id int) (ObservationMetadata, error)
	CreateObservationMetadata(observationId int, long float64, lat float64) (int, error)
}
