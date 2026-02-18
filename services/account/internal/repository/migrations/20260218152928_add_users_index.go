package migrations

import (
	"context"
	"database/sql"
)

func upAddUsersIndexGo(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_user_login;`)
	return err
}

func downAddUsersIndexGo(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DRFOP INDEX IF EXISTS idx_user_login`)
	return err
}
