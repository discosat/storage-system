package observationMetadata

import "time"

type ObservationMetadata struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Long      int       `json:"long"`
	Lat       int       `json:"lat"`
}
