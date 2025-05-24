package observationRequest

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/discosat/storage-system/internal/Commands"
	"log/slog"
)

type PsqlObservationRequestRepository struct {
	db *sql.DB
}

func NewPsqlObservationRequestRepository(db *sql.DB) ObservationRequestRepository {
	return &PsqlObservationRequestRepository{db: db}
}

func (p PsqlObservationRequestRepository) GetObservationRequest(id int) (ObservationRequestAggregate, error) {

	var observationRequestEntity ObservationRequestAggregate
	err := p.db.QueryRow("SELECT o_r.id, o_r.type, fp.id, fp.name, fp.user_id, m.id, m.name, m.bucket FROM observation_request o_r INNER JOIN public.flight_plan fp ON o_r.flight_plan_id = fp.id INNER JOIN public.mission m ON m.id = fp.mission_id WHERE o_r.id = $1", id).
		Scan(&observationRequestEntity.ObservationRequest.Id,
			//&observationRequestEntity.ObservationRequest.CreatedAt,
			//&observationRequestEntity.ObservationRequest.UpdatedAt,
			//&observationRequestEntity.ObservationRequest.FlightPlanId,
			&observationRequestEntity.ObservationRequest.OType,
			&observationRequestEntity.FlightPlan.Id,
			//&observationRequestEntity.FlightPlanEntity.CreatedAt,
			//&observationRequestEntity.FlightPlanEntity.UpdatedAt,
			&observationRequestEntity.FlightPlan.Name,
			&observationRequestEntity.FlightPlan.UserId,
			//&observationRequestEntity.FlightPlanEntity.MissionId,
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

func (p PsqlObservationRequestRepository) GetFlightPlanById(id int) (FlightPlanAggregate, error) {
	var flightPlan FlightPlanAggregate
	tx, err := p.db.BeginTx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return FlightPlanAggregate{}, err
	}
	err = tx.QueryRow("SELECT id, name, user_id, mission_id, locked FROM flight_plan WHERE id = $1", id).Scan(&flightPlan.Id, &flightPlan.Name, &flightPlan.UserId, &flightPlan.MissionId, &flightPlan.Locked)
	if err != nil {
		return FlightPlanAggregate{}, err
	}

	rows, err := tx.Query("SELECT id, type FROM observation_request where flight_plan_id = $1", flightPlan.Id)
	if err != nil {
		return FlightPlanAggregate{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var observationRequest ObservationRequestDTO
		if err := rows.Scan(&observationRequest.Id, &observationRequest.OType); err != nil {
			return FlightPlanAggregate{}, err
		}
		flightPlan.ObservationRequests = append(flightPlan.ObservationRequests, observationRequest)
	}
	tx.Commit()
	return flightPlan, nil
}

func (p PsqlObservationRequestRepository) CreateFlightPlan(flightPlan Commands.CreateFlightPlanCommand, requestList []Commands.CreateObservationRequestCommand) (int, error) {

	slog.Info(fmt.Sprintf("Creating a flightplan: %v, for missionId: %v, with observation requests: %v", flightPlan.Name, flightPlan.MissionId, requestList))
	tx, err := p.db.BeginTx(context.Background(), &sql.TxOptions{})
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
			return -1, &ObservationRequestError{msg: "Observation request is formatted wrong", code: ObservationRequestParseError}
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

func (p PsqlObservationRequestRepository) UpdateFlightPlan(flightPlan FlightPlanAggregate) (int, error) {
	tx, err := p.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return -1, err
	}
	var locked bool
	err = tx.QueryRow("SELECT locked FROM flight_plan WHERE id = $1", flightPlan.Id).Scan(&locked)
	if err != nil {
		return -1, err
	}
	if locked {
		return flightPlan.Id, &ObservationRequestError{msg: fmt.Sprintf("Flight plan with id: %v is locked", flightPlan.Id), code: FlightPlanIsLocked}
	}

	_, err = tx.Exec("UPDATE flight_plan SET name = $1 WHERE id = $2", flightPlan.Name, flightPlan.Id)
	if err != nil {
		return -1, err
	}

	deleteIds := make(map[int]bool)
	rows, err := tx.Query("SELECT id FROM observation_request WHERE flight_plan_id = $1", flightPlan.Id)
	if err != nil {
		return -1, err
	}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return -1, err
		}
		deleteIds[id] = true
	}

	for _, request := range flightPlan.ObservationRequests {
		deleteIds[request.Id] = false
	}

	for _, request := range flightPlan.ObservationRequests {
		if deleteIds[request.Id] {
			_, err := tx.Exec("DELETE FROM observation_request WHERE id = $1", request.Id)
			if err != nil {
				return -1, err
			}
		} else {
			_, err := tx.Exec("UPDATE observation_request SET type = $1 WHERE id = $2", request.OType, request.Id)
			if err != nil {
				return 0, err
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return flightPlan.Id, nil
}
func (p PsqlObservationRequestRepository) DeleteFlightPlan(id int) (bool, error) {
	return false, nil
}

func (p PsqlObservationRequestRepository) GetObservationRequestById(id int) (ObservationRequest, error) {
	var observationRequest ObservationRequest
	err := p.db.QueryRow("SELECT * FROM observation_request WHERE id = $1", id).Scan(&observationRequest.Id, &observationRequest.FlightPlanId, &observationRequest.OType)
	return observationRequest, err
}

func (p PsqlObservationRequestRepository) CreateObservationRequest(flightPLanId int, camera string) {
	//TODO implement me
	panic("implement me")
}
