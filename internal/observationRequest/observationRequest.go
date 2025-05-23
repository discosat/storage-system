package observationRequest

import "time"

type ObservationRequest struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	OType     string    `json:"o_type"`
	FId       int       `json:"f_id"`
}

type ObservationRequestAggregate struct {
	FlightPlan         FlightPlanDTO         `json:"flightPlan"`
	Mission            MissionDTO            `json:"mission"`
	ObservationRequest ObservationRequestDTO `json:"observationRequest"`
}

type ObservationRequestDTO struct {
	Id    int    `json:"id"`
	OType string `json:"o_type"`
}
