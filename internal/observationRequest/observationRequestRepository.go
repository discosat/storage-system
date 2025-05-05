package observationRequest

type ObservationRequestRepository interface {
	GetObservationRequest(id int) (ObservationRequestAggregate, error)
	GetObservationRequestById(id int) (ObservationRequest, error)
	GetFlightPlantById(id int) (FlightPlan, error)
	GetMissionById(id int) (Mission, error)
}
