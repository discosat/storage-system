package observationRequest

type ObservationRequestRepository interface {
	GetObservationRequest(id int) (ObservationRequestCommand, error)
	GetFlightPlanById(id int) (FlightPlanAggregate, error)
	UpdateFlightPlan(flightplan FlightPlanAggregate) (int, error)
	DeleteFlightPlan(id int) (bool, error)
	CreateFlightPlan(flightPlan FlightPlanAggregate) (int, error)
}
