package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
    Host     string
    Port     string
    LogLevel string
    LogFile  string
}

func loggingMiddleware(h http.Handler) http.Handler {
    logFn := func(rw http.ResponseWriter, r *http.Request) {
        start := time.Now()

        uri := r.RequestURI
        method := r.Method
        h.ServeHTTP(rw, r) // serve the original request

        duration := time.Since(start)

        // log request details
        log.WithFields(log.Fields{
            "uri":      uri,
            "method":   method,
            "duration": duration,
        }).Info("Request handled")
    }
    return http.HandlerFunc(logFn)
}

func addRoutes(mux *http.ServeMux, logger *log.Logger, config *Config) {
    mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, world!")
    }))
}

func NewServer(logger *log.Logger, config *Config) http.Handler {
    mux := http.NewServeMux()
    addRoutes(mux, logger, config)

    handler := loggingMiddleware(mux)
    return handler
}

func main() {
    config := &Config{
        Host:     "localhost",
        Port:     "8080",
        LogLevel: "info",
        LogFile:  "server.log",
    }
    logger := log.New()

    srv := NewServer(logger, config)
    httpServer := &http.Server{
        Addr:    net.JoinHostPort(config.Host, config.Port),
        Handler: srv,
    }

    go func() {
        log.Printf("Listening on %s\n", httpServer.Addr)
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Fprintf(os.Stderr, "Error listening and serving: %s\n", err)
        }
    }()

    // Wait for interrupt signal to gracefully shut down the server with a timeout
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)

    <-stop // Wait for SIGINT (Ctrl+C)
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := httpServer.Shutdown(ctx); err != nil {
        log.Printf("Error shutting down http server: %s\n", err)
    }

    log.Println("Server gracefully stopped")
}
