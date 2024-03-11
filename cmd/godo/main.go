package main

import (
	"log/slog"
	"os"

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

	srv := server.New(logger, config)
	server.Run(logger, srv)
}
