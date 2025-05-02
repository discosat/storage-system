package observation

import "mime/multipart"

type ObservationRepository interface {
	GetById(id int) (Observation, error)
	// CreateObservation TODO Should preferably not depend on specific file type. Check if a Reader could do
	CreateObservation(file *multipart.FileHeader, bucket string, flightPlanName string, observationRequestId int) (int, error)
}
