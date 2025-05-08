package observationRequest

type ObservationRequestRepository interface {
	GetObservationRequest(id int) (ObservationRequestAggregate, error)
	GetObservationRequestById(id int) (ObservationRequest, error)
	CreateObservationRequest(flightPLanId int, camera string)
	GetFlightPlantById(id int) (FlightPlan, error)
	CreateFlightPlan(missionId int, userId int, name string, requestList []ObservationRequest) (int, error)
	GetMissionById(id int) (Mission, error)
	//CreateMission() ()
}
