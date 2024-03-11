package main

import (
	"log/slog"

	"github.com/confy/godo/internal/server"
)



func main() {
	config, err := server.LoadConfig()
	if err != nil {
		panic(err)
	}
	logger := server.GetLogger(config)
	slog.SetDefault(logger)

	srv := server.New(logger, config)
	server.Run(logger, srv)

}
