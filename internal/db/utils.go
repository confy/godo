package db

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var Schema string

func CreateTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, Schema)
	return err
}
