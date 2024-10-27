package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddSomeColumn, downAddSomeColumn)
}

func upAddSomeColumn(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			name TEXT,
			surname TEXT);`)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO users VALUES
			(0, 'root', '', ''),
			(1, 'vojtechvitek', 'Vojtech', 'Vitek');`)

	return err
}

func downAddSomeColumn(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE users;")
	return err
}
