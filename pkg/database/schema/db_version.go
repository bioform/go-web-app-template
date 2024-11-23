package schema

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func DBVersion(db *sql.DB) (int64, error) {
	// Get the most recent migration version
	version, err := goose.GetDBVersion(db)

	return version, err
}
