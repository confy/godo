package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/a-h/templ"
	"github.com/confy/godo/middleware"
	"github.com/confy/godo/views"
	// "github.com/go-kit/log"
)

type Config struct {
	Host     string
	Port     string
	LogLevel slog.Level
}

func loadConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     "8080",
		LogLevel: slog.LevelDebug,
	}
}

func NewServer(logger *slog.Logger, config *Config) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)

	handler := middleware.LoggingMiddleware(logger)(mux)
	return handler
}

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handleRoot)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templ.Handler(views.IndexPage("Hello, world!")).ServeHTTP(w, r)
}

func startServer(logger *slog.Logger, httpServer *http.Server) {
	go func() {
		logger.Info("Starting server...", "host", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("msg", "Error listening and serving", "err", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // Wait for SIGINT (Ctrl+C)
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("msg", "Error shutting down http server", "err", err)
	}

	logger.Info("Server gracefully stopped")
}

func main() {
	config := loadConfig()
	logOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logHandler := slog.NewTextHandler(os.Stderr, logOpts)
	// logger := slog.
	logger := slog.New(logHandler)
	srv := NewServer(logger, config)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	startServer(logger, httpServer)
}
