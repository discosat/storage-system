package observation

type ObservationRepository interface {
	GetById(id int) (Observation, error)
	GetByRequest(requestId int) (Observation, error)
	CreateObservation(requestId int, userId int) (Observation, error)
}
