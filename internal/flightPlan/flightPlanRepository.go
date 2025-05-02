package flightPlan

type FlightPlanRepository interface {
	GetById(id string) (FlightPlan, error)
}
