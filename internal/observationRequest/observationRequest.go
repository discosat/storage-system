package observationRequest

import "time"

type ObservationRequest struct {
	Id           int       `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	OType        string    `json:"o_type"`
	FlightPlanId int       `json:"flight_plan_id"`
}

type ObservationRequestCommand struct {
	FlightPlanName     string
	Bucket             string
	ObservationRequest ObservationRequestDTO `json:"observationRequest"`
}

type ObservationRequestDTO struct {
	Id    int    `json:"id"`
	OType string `json:"o_type"`
}
