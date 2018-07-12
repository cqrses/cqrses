package mysql

import (
	"database/sql"
)

type (
	EventStore struct {
		db *sql.DB
	}
)
