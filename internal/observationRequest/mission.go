package observationRequest

import "time"

type Mission struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Bucket    string    `json:"bucket"`
}

type MissionDTO struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Bucket string `json:"bucket"`
}
