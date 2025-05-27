package observationRequest

type FlightPlanAggregate struct {
	Id                  int                     `json:"id"`
	Name                string                  `json:"name"`
	UserId              int                     `json:"user_id"`
	MissionId           int                     `json:"mission_id"`
	Locked              bool                    `json:"locked"`
	ObservationRequests []ObservationRequestDTO `json:"observation_requests"`
}

type FlightPlanDTO struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	UserId    int    `json:"user_id"`
	MissionId int    `json:"mission_id"`
	Locked    bool   `json:"locked"`
}
