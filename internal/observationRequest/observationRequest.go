package observationRequest

import "time"

type ObservationRequest struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	OType     string    `json:"o_type"`
	RId       int       `json:"r_id"`
}
