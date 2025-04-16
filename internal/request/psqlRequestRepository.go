package request

import "database/sql"

type PsqlRequestRepository struct {
	db *sql.DB
}

func NewPsqlRequestRepository(db *sql.DB) RequestRepository {
	return &PsqlRequestRepository{db: db}
}

func (p PsqlRequestRepository) GetById(id string) (Request, error) {
	var request Request
	err := p.db.QueryRow("SELECT * FROM request WHERE id = $1", id).Scan(&request.Id, &request.Name, &request.UserId, &request.MissionId)
	if err != nil {
		return request, err
	}
	return request, nil
}
