package observationRequest

type ObservationRequestRepository interface {
	GetById(id int) (ObservationRequest, error)
}
