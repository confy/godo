package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/confy/godo/internal/db"
	"github.com/confy/godo/internal/server"
)

func main() {
	config, err := server.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     config.LogLevel,
		}),
	)
	slog.SetDefault(logger)
	dbUrl, err := config.GetDbURL()
	if err != nil {
		panic(err)
	}

	// connect to the database
	conn, err := sql.Open("libsql", dbUrl)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ctx := context.Background()

	// create tables
	err = db.CreateTables(ctx, conn)
	if err != nil {
		panic(err)
	}

	database := db.New(conn)
	database.CreateUser(ctx, db.CreateUserParams{
		Email:    "me@adrian.ooo",
		Username: "Adrian",
	})

	srv := server.New(logger, config, database)
	server.Run(logger, srv)
}
