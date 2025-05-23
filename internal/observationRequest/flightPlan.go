package observationRequest

import "time"

type FlightPlanEntity struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	UserId    int       `json:"user_id"`
	MissionId int       `json:"mission_id"`
}

type FlightPlanDTO struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
}
