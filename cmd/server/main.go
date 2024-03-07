package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/a-h/templ"
	"github.com/confy/godo/middleware"
	"github.com/confy/godo/views"
	"github.com/go-kit/log"
)

type Config struct {
	Host string
	Port string
}

func loadConfig() *Config {
	return &Config{
		Host: "localhost",
		Port: "8080",
	}
}

func NewServer(logger log.Logger, config *Config) http.Handler {
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

func startServer(logger log.Logger, httpServer *http.Server) {
	go func() {
		logger.Log("msg", "Starting server...", "host", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log("err", fmt.Sprintf("Error listening and serving: %s", err))
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
		logger.Log("err", fmt.Sprintf("Error shutting down http server: %s", err))
	}

	logger.Log("msg", "Server gracefully stopped")
}

func main() {
	config := loadConfig()

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	srv := NewServer(logger, config)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	startServer(logger, httpServer)
}
