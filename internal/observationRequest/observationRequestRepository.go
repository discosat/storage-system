package observationRequest

type ObservationRequestRepository interface {
	GetObservationRequest(id int) (ObservationRequestAggregate, error)
	GetObservationRequestById(id int) (ObservationRequest, error)
	CreateObservationRequest(flightPLanId int, camera string)
	GetFlightPlanById(id int) (FlightPlanAggregate, error)
	UpdateFlightPlan(flightplan FlightPlanAggregate) (int, error)
	DeleteFlightPlan(id int) (bool, error)
	CreateFlightPlan(flightPlan FlightPlanAggregate) (int, error)
	GetMissionById(id int) (Mission, error)
}
