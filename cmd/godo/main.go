package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	middleware "github.com/confy/godo/middleware"
	log "github.com/go-kit/log"
)

type Config struct {
	Host     string
	Port     string
	LogLevel string
	LogFile  string
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Hello, world!")
	}))
}

func NewServer(logger log.Logger, config *Config) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)

	handler := middleware.LoggingMiddleware(logger)(mux)
	return handler
}

func main() {
	config := &Config{
		Host: "localhost",
		Port: "8080",
	}
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	srv := NewServer(logger, config)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	go func() {
		logger.Log("msg", "Starting server...", "host", config.Host, "port", config.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Error listening and serving: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // Wait for SIGINT (Ctrl+C)
	logger.Log("msg", "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Log("err", "Error shutting down http server: %s", err)
	}

	logger.Log("msg", "Server gracefully stopped")
}
