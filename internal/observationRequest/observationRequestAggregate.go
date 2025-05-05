package observationRequest

type ObservationRequestAggregate struct {
	FlightPlan         FlightPlanDTO         `json:"flightPlan"`
	Mission            MissionDTO            `json:"mission"`
	ObservationRequest ObservationRequestDTO `json:"observationRequest"`
}
