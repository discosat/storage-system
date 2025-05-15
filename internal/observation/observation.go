package observation

import "time"

type Observation struct {
	Id                   int       `json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	ObjectReference      string    `json:"objectReference"`
	ObservationRequestId int       `json:"observationRequestId"`
}
