package request

type Request struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	UserId    int    `json:"user_id"`
	MissionId int    `json:"mission_id"`
}
