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

	// create a session manager
	sessionManager := server.CreateSessionManager(config.UseHttps)

	// connect to the database
	dbUrl := config.GetDbURL()
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

	srv := server.New(logger, config, sessionManager, database)
	server.Run(logger, srv)
}
