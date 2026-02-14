package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddNamedMigrationContext("20260214102606_initdb.go", upInitdb, downInitdb)
}

func upInitdb(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS users(
				id BIGSERIAL PRIMARY KEY,
				login text NOT NULL UNIQUE,
				email TEXT NOT NULL UNIQUE,
				phone TEXT,
				first_name TEXT,
				last_name TEXT,
				middle_name TEXT,
				age INT,
				created_at TIMESTAMP NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMP NOT NULL DEFAULT NOW()
			);
		`)
	return err
}

func downInitdb(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS users;`)
	return err
}
