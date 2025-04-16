package measurementRequest

type MeasurementRequestRepository interface {
	GetById(id int) (MeasurementRequest, error)
}
