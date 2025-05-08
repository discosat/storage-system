package observationRequest

type ObservationRequestRepository interface {
	GetObservationRequest(id int) (ObservationRequestAggregate, error)
	GetObservationRequestById(id int) (ObservationRequest, error)
	CreateObservationRequest(flightPLanId int, camera string)
	GetFlightPlantById(id int) (FlightPlan, error)
	CreateFlightPlan(flightPlan FlightPlanCommand, requestList []ObservationRequestCommand) (int, error)
	GetMissionById(id int) (Mission, error)
	//CreateMission() ()
}
