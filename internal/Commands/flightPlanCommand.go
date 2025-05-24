package Commands

type CreateFlightPlanCommand struct {
	Name      string `json:"name"`
	UserId    int    `json:"userId"`
	MissionId int    `json:"missionId"`
}
