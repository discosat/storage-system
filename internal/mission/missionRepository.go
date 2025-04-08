package mission

import "database/sql"

type MissionRepository interface {
	GetById(db *sql.DB, id int) (Mission, error)
}
