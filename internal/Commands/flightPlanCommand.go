package Commands

type FlightPlanCommand struct {
	Name      string `json:"name"`
	UserId    int    `json:"userId"`
	MissionId int    `json:"missionId"`
}
