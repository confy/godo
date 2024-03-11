package server

import (
	"log/slog"
	"os"
)

func GetLogger(config *Config) *slog.Logger {
	logOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     config.LogLevel,
	}
	logHandler := slog.NewTextHandler(os.Stderr, logOpts)
	logger := slog.New(logHandler)
	return logger
}