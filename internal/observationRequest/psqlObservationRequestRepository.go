package observationRequest

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/discosat/storage-system/internal/Commands"
	"github.com/jmoiron/sqlx"
	"log/slog"
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

func (p PsqlObservationRequestRepository) CreateFlightPlan(flightPlan Commands.FlightPlanCommand, requestList []Commands.ObservationRequestCommand) (int, error) {

	slog.Info(fmt.Sprintf("Creating a flightplan: %v, for missionId: %v, with observation requests: %v", flightPlan.Name, flightPlan.MissionId, requestList))
	tx, err := p.db.BeginTxx(context.Background(), &sql.TxOptions{})
	defer tx.Rollback()
	if err != nil {
		slog.Error(fmt.Sprintf("Could not create transaction: %v", err))
		return -1, err
	}

	var fpId int
	rows, err := tx.Query("INSERT INTO flight_plan (name, user_id, mission_id) VALUES ($1, $2, $3) RETURNING id", flightPlan.Name, flightPlan.UserId, flightPlan.MissionId)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not insert flightplan: %v", err))
		return -1, err
	}

	rows.Next()
	err = rows.Scan(&fpId)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not bind flightPlant id: %v", err))
		return -1, err
	}

	err = rows.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("Flight plan row could not be closed: %v", err))
		return -1, err
	}

	for _, request := range requestList {
		_, qErr := tx.Exec("INSERT INTO observation_request (flight_plan_id, type) VALUES ($1, $2)", fpId, request.OType)
		if qErr != nil {
			slog.Error(fmt.Sprintf("Formatting eror of observation request: %v. \n Error: %v", request, err))
			return -1, err
		}
	}
	err = tx.Commit()
	if err != nil {
		slog.Error(fmt.Sprintf("Could not commit transaction: %v", err))
		return -1, err
	}
	slog.Info(fmt.Sprintf("Flight plan: %v, with id %v, created", flightPlan.Name, fpId))
	return fpId, nil
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
