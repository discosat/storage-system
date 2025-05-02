package flightPlan

import "time"

type FlightPlan struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	UserId    int       `json:"user_id"`
	MissionId int       `json:"mission_id"`
}
