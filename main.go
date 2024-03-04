package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)


type Config struct {
    Host string
    Port string
    LogLevel string
    LogFile string
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
        })
    }
    return http.HandlerFunc(logFn)
}


func addRoutes(
	mux                 *http.ServeMux,
	logger              *log.Logger,
	config              *Config,
) {
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, world!")
    }))
}

func NewServer(
	logger *log.Logger,
	config *Config,

) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		config,
	)
	var handler http.Handler = mux
    handler = loggingMiddleware(handler)
	return handler
}

func main() {
    config := &Config{
        Host: "localhost",
        Port: "8080",
        LogLevel: "info",
        LogFile: "server.log",
    }
    logger := log.New()

    srv := NewServer(
        logger,
        config,
    )
    httpServer := &http.Server{
        Addr:    net.JoinHostPort(config.Host, config.Port),
        Handler: srv,
    }
    go func() {
        log.Printf("listening on %s\n", httpServer.Addr)
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
        }
    }()
    var wg sync.WaitGroup
    wg.Add(1)

    ctx, _ := context.WithCancel(context.Background())
    go func() {
        defer wg.Done()
        <-ctx.Done()
        // make a new context for the Shutdown (thanks Alessandro Rosetti)
        shutdownCtx := context.Background()
        shutdownCtx, cancel := context.WithTimeout(ctx, 10 * time.Second)
        defer cancel()
        if err := httpServer.Shutdown(shutdownCtx); err != nil {
            fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
        }
    }()
    wg.Wait()
    log.Println("server stopped")
}