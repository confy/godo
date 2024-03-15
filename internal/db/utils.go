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

// Create a new user or get an existing user
func CreateOrGetUser(ctx context.Context, dbQueries *Queries, user CreateUserParams) (User, error) {
	dbUser, err := dbQueries.GetUserByLogin(ctx, user.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return dbQueries.CreateUser(ctx, user)
		}
	}
	return dbUser, nil
}
