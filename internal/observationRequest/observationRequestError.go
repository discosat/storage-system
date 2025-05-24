package observationRequest

const (
	_ = iota
	FlightPlanParseError
	ObservationRequestParseError
	FlightPlanIsLocked
)

type ObservationRequestError struct {
	msg  string
	code int
}

func (e *ObservationRequestError) Error() string { return e.msg }

func (e *ObservationRequestError) Code() int { return e.code }
