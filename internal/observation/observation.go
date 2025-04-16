package observation

type Observation struct {
	Id        int `json:"id"`
	RequestId int `json:"request_id"`
	UserId    int `json:"user_id"`
}
