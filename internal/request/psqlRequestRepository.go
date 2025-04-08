package request

import "database/sql"

type PsqlRequestRepository struct {
	RequestRepository RequestRepository
}

func (p PsqlRequestRepository) GetById(db *sql.DB, id string) (Request, error) {
	var request Request
	err := db.QueryRow("SELECT * FROM request WHERE id = $1", id).Scan(&request.Id, &request.Name, &request.UserId, &request.MissionId)
	if err != nil {
		return request, err
	}
	return request, nil
}
