package Commands

type CreateObservationRequestCommand struct {
	OType string `json:"o_type"`
}
type UpdateObservationRequestCommand struct {
	Id    int    `json:"id"`
	OType string `json:"o_type"`
}
