package mission

import (
	"database/sql"
)

type PsqlMissionRepository struct {
	MissionRepository MissionRepository
}

func (p PsqlMissionRepository) GetById(db *sql.DB, id int) (Mission, error) {
	var mission Mission
	err := db.QueryRow("SELECT * FROM mission WHERE id = $1", id).Scan(&mission.Id, &mission.Name, &mission.Bucket)
	if err != nil {
		return mission, err
	}
	return mission, nil
}
