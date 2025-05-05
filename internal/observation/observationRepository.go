package observation

import "mime/multipart"

type ObservationRepository interface {
	GetObservation(id int) (Observation, error)
	// CreateObservation TODO Should preferably not depend on specific file type. Check if a Reader could do
	CreateObservation(file *multipart.FileHeader, bucket string, flightPlanName string, observationRequestId int) (int, error)

	GetObservationMetadata(id int) (ObservationMetadata, error)
	CreateObservationMetadata(observationId int, long float64, lat float64) (int, error)
}
