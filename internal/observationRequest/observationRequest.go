package observationRequest

import "time"

type ObservationRequest struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	OType     string    `json:"o_type"`
	FId       int       `json:"f_id"`
}
