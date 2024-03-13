package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/confy/godo/internal/db"
	"github.com/confy/godo/internal/middleware"
	"github.com/confy/godo/views"
)

func New(logger *slog.Logger, config *Config, sessionManager *scs.SessionManager, db *db.Queries) *http.Server {
	mux := http.NewServeMux()
	addRoutes(mux, db)

	// var githubOauthConfig = &oauth2.Config{
	// 	RedirectURL:  config.GetHostUrl() + "/callback",
	// 	ClientID:     config.GithubClientID,
	// 	ClientSecret: config.GithubClientSecret,
	// 	Endpoint:     github.Endpoint,
	// }

	handler := middleware.LoggingMiddleware(logger)(mux)
	handler = sessionManager.LoadAndSave(handler)
	server := &http.Server{
		Addr:     net.JoinHostPort(config.Host, config.Port),
		Handler:  handler,
		ErrorLog: slog.NewLogLogger(logger.Handler(), config.LogLevel),
	}
	return server
}

func addRoutes(mux *http.ServeMux, db *db.Queries) {
	mux.HandleFunc("/", handleRoot)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templ.Handler(views.IndexPage("Hello, world!")).ServeHTTP(w, r)
}

func Run(logger *slog.Logger, httpServer *http.Server) {
	go func() {
		logger.Info("Starting server...", "host", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error listening and serving", "err", err)
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
		logger.Error("Error shutting down http server", "err", err)
	}

	logger.Info("Server gracefully stopped")
}
