package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/confy/godo/config"
	"github.com/confy/godo/internal/db"
	"github.com/confy/godo/internal/libsqlstore"
	"github.com/confy/godo/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     cfg.LogLevel,
		}),
	)
	slog.SetDefault(logger)

	// Connect to the database
	dbURL := cfg.GetDBURL()
	conn, err := sql.Open("libsql", dbURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ctx := context.Background()

	// Create tables
	err = db.CreateTables(ctx, conn)
	if err != nil {
		panic(err)
	}

	database := db.New(conn)

	// Create a new session manager - we should change to a new session db eventually
	session := scs.New()
	session.Store = libsqlstore.New(conn)
	session.Cookie.SameSite = http.SameSiteStrictMode
	session.Cookie.Secure = cfg.UseHTTPS

	srv := server.New(cfg, database, session)
	server.Run(srv)
}
