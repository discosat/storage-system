package observationRequest

import "time"

type Mission struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Bucket    string    `json:"bucket"`
}
