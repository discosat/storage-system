package mission

import (
	"database/sql"
)

type PsqlMissionRepository struct {
	db *sql.DB
}

func NewPsqlMissionRepository(db *sql.DB) MissionRepository {
	return &PsqlMissionRepository{db: db}
}

func (p PsqlMissionRepository) GetById(id int) (Mission, error) {
	var mission Mission
	err := p.db.QueryRow("SELECT * FROM mission WHERE id = $1", id).Scan(&mission.Id, &mission.CreatedAt, &mission.UpdatedAt, &mission.Name, &mission.Bucket)
	if err != nil {
		return mission, err
	}
	return mission, nil
}
