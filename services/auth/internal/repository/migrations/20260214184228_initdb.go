package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitdb, downInitdb)
}

func upInitdb(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS users(
				id BIGSERIAL PRIMARY KEY,
				login TEXT NOT NULL UNIQUE,
				email TEXT NOT NULL UNIQUE,
				password_hash NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMP NOT NULL DEFAULT NOW()
			);
		`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS refresh_token(
				id BIGSERIAL PRIMARY KEY,
				user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				token TEXT NOT NULL UNIQUE,
				expires_at TIMESTAMP NOT NULL DEFAULT NOW(),
				revoked_at TIMESTAMP NULL,
				created_at TIMESTAMP NOT NULL DEFAULT NOW()
			);
		`)
	if err != nil {
		return err
	}
	return nil
}

func downInitdb(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS refresh_tokens;`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DROP TABLE IF EXISTS users;`)
	if err != nil {
		return err
	}

	return nil
}
