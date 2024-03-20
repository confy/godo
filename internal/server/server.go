package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/confy/godo/internal/db"
	"github.com/confy/godo/internal/handler"
	"github.com/confy/godo/internal/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func addRoutes(
	mux *http.ServeMux,
	session *scs.SessionManager,
	oauth *oauth2.Config,
	database *db.Queries,
) {
	mux.HandleFunc("/login", handler.HandleAuthLogin(oauth))
	mux.HandleFunc("/callback", handler.HandleAuthCallback(session, oauth, database))
	mux.HandleFunc("/", handler.HandleRoot(database, session))
	mux.HandleFunc("/test", middleware.RequireLogin(handler.HandleTestPage(database, session), session))
}

func New(config *Config, database *db.Queries, session *scs.SessionManager) *http.Server {
	mux := http.NewServeMux()

	oauth := &oauth2.Config{
		RedirectURL:  config.GetHostURL() + "/callback",
		ClientID:     config.GithubClientID,
		ClientSecret: config.GithubClientSecret,
		Endpoint:     github.Endpoint,
	}

	addRoutes(mux, session, oauth, database)

	handler := middleware.LoggingMiddleware()(mux)
	handler = session.LoadAndSave(handler)

	server := &http.Server{
		Addr:              net.JoinHostPort(config.Host, config.Port),
		Handler:           handler,
		ErrorLog:          slog.NewLogLogger(slog.Default().Handler(), config.LogLevel),
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	return server
}

func Run(httpServer *http.Server) {
	go func() {
		slog.Info("Starting server...", "host", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error listening and serving", "err", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // Wait for SIGINT (Ctrl+C)
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down http server", "err", err)
	}

	slog.Info("Server gracefully stopped")
}
