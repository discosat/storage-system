package observationRequest

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
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
	var flightPlan FlightPlan
	err := p.db.QueryRow("SELECT * FROM flight_plan WHERE id = $1", id).Scan(&flightPlan.Id, &flightPlan.CreatedAt, &flightPlan.UpdatedAt, &flightPlan.Name, &flightPlan.UserId, &flightPlan.MissionId)
	if err != nil {
		return flightPlan, err
	}
	return flightPlan, nil
}

func (p PsqlObservationRequestRepository) CreateFlightPlan(missionId int, userId int, name string, requestList []ObservationRequest) (int, error) {

	tx, err := p.db.BeginTxx(context.Background(), &sql.TxOptions{})
	//defer tx.Rollback()
	if err != nil {
		return -1, err
	}

	var fpId int
	rows, err := tx.Query("INSERT INTO flight_plan (name, user_id, mission_id) VALUES ($1, $2, $3) RETURNING id", name, userId, missionId)
	rows.Next()
	rows.Scan(&fpId)
	if err != nil {
		return -1, err
	}
	rows.Close()
	for _, request := range requestList {
		r, qErr := tx.Query("INSERT INTO observation_request (flight_plan_id, type) VALUES ($1, $2)", fpId, request.OType)
		if qErr != nil {
			log.Fatalf("%v", qErr)
		}
		r.Scan()
		r.Close()
		//log.Println(r)
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return fpId, nil

	//TODO implement me
	panic("implement me")
}

func (p PsqlObservationRequestRepository) GetObservationRequestById(id int) (ObservationRequest, error) {
	var observationRequest ObservationRequest
	err := p.db.QueryRow("SELECT * FROM observation_request WHERE id = $1", id).Scan(&observationRequest.Id, &observationRequest.FId, &observationRequest.OType)
	return observationRequest, err
}

func (p PsqlObservationRequestRepository) CreateObservationRequest(flightPLanId int, camera string) {
	//TODO implement me
	panic("implement me")
}
