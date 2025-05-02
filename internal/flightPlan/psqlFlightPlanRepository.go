package flightPlan

import "database/sql"

type PsqlFlightPlanRepository struct {
	db *sql.DB
}

func NewPsqlFlightPlanRepository(db *sql.DB) FlightPlanRepository {
	return &PsqlFlightPlanRepository{db: db}
}

func (p PsqlFlightPlanRepository) GetById(id string) (FlightPlan, error) {
	var flightPlan FlightPlan
	err := p.db.QueryRow("SELECT * FROM flight_plan WHERE id = $1", id).Scan(&flightPlan.Id, &flightPlan.CreatedAt, &flightPlan.UpdatedAt, &flightPlan.Name, &flightPlan.UserId, &flightPlan.MissionId)
	if err != nil {
		return flightPlan, err
	}
	return flightPlan, nil
}
