package observationRequest

import "github.com/discosat/storage-system/internal/Commands"

type ObservationRequestRepository interface {
	GetObservationRequest(id int) (ObservationRequestAggregate, error)
	GetObservationRequestById(id int) (ObservationRequest, error)
	CreateObservationRequest(flightPLanId int, camera string)
	GetFlightPlantById(id int) (FlightPlanEntity, error)
	CreateFlightPlan(flightPlan Commands.FlightPlanCommand, requestList []Commands.ObservationRequestCommand) (int, error)
	GetMissionById(id int) (Mission, error)
}
