package Commands

import (
	"bytes"
)

type CreateObservationCommand struct {
	File                 *bytes.Reader
	FileSize             int64
	FileName             string
	Bucket               string
	FlightPlanName       string
	ObservationRequestId int
	UserId               int
	ObservationType      string
}
