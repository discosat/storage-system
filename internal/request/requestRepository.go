package request

import "database/sql"

type RequestRepository interface {
	GetById(db *sql.DB, id string) (Request, error)
}
