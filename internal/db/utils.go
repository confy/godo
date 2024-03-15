package db

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
)

//go:embed schema.sql
var Schema string

func dropTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "DROP TABLE IF EXISTS users")
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "DROP TABLE IF EXISTS todos")
	return err
}

func CreateTables(ctx context.Context, db *sql.DB) error {
	err := dropTables(ctx, db) // delete tables during development

	if err != nil {
		panic(err)
	}
	_, err = db.ExecContext(ctx, Schema)
	return err
}

// Create a new user or get an existing user
func CreateOrGetUser(ctx context.Context, dbQueries *Queries, user CreateUserParams) (User, error) {
	dbUser, err := dbQueries.GetUserByLogin(ctx, user.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dbQueries.CreateUser(ctx, user)
		}
	}
	return dbUser, nil
}
