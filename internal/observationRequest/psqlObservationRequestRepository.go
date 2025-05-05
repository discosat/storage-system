package observationRequest

import (
	"github.com/jmoiron/sqlx"
)

type PsqlObservationRequestRepository struct {
	db *sqlx.DB
}

func NewPsqlObservationRequestRepository(db *sqlx.DB) ObservationRequestRepository {
	return &PsqlObservationRequestRepository{db: db}
}

func (p PsqlObservationRequestRepository) GetObservationRequest(id int) (ObservationRequestAggregate, error) {

	var observationRequestEntity ObservationRequestAggregate
	err := p.db.QueryRowx("SELECT o_r.id, o_r.type, fp.id, fp.name, fp.user_id, m.id, m.name, m.bucket FROM observation_request o_r INNER JOIN public.flight_plan fp ON o_r.flight_plan_id = fp.id INNER JOIN public.mission m ON m.id = fp.mission_id WHERE o_r.id = $1", id).
		Scan(&observationRequestEntity.ObservationRequest.Id,
			//&observationRequestEntity.ObservationRequest.CreatedAt,
			//&observationRequestEntity.ObservationRequest.UpdatedAt,
			//&observationRequestEntity.ObservationRequest.FId,
			&observationRequestEntity.ObservationRequest.OType,
			&observationRequestEntity.FlightPlan.Id,
			//&observationRequestEntity.FlightPlan.CreatedAt,
			//&observationRequestEntity.FlightPlan.UpdatedAt,
			&observationRequestEntity.FlightPlan.Name,
			&observationRequestEntity.FlightPlan.UserId,
			//&observationRequestEntity.FlightPlan.MissionId,
			&observationRequestEntity.Mission.Id,
			//&observationRequestEntity.Mission.CreatedAt,
			//&observationRequestEntity.Mission.UpdatedAt,
			&observationRequestEntity.Mission.Name,
			&observationRequestEntity.Mission.Bucket)

	return observationRequestEntity, err
}

func (p PsqlObservationRequestRepository) GetMissionById(id int) (Mission, error) {
	//TODO implement me
	panic("implement me")
}

func (p PsqlObservationRequestRepository) GetFlightPlantById(id int) (FlightPlan, error) {
	//TODO implement me
	panic("implement me")
}

func (p PsqlObservationRequestRepository) GetObservationRequestById(id int) (ObservationRequest, error) {
	var observationRequest ObservationRequest
	err := p.db.QueryRow("SELECT * FROM observation_request WHERE id = $1", id).Scan(&observationRequest.Id, &observationRequest.FId, &observationRequest.OType)
	return observationRequest, err
}
