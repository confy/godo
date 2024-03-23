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
	"github.com/confy/godo/config"
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
	mux.HandleFunc("GET /login", handler.HandleAuthLogin(oauth))
	mux.HandleFunc("GET /logout", handler.HandleAuthLogout(session))
	mux.HandleFunc("GET /callback", handler.HandleAuthCallback(session, oauth, database))

	mux.HandleFunc("GET /", handler.HandleRoot(database, session))
	mux.HandleFunc("GET /todos", middleware.RequireLogin(handler.HandleTodos(database, session), session))

	mux.HandleFunc("DELETE /todo/{id}", middleware.RequireLogin(handler.HandleDeleteTodo(database, session), session))
}

func New(cfg *config.Config, database *db.Queries, session *scs.SessionManager) *http.Server {
	mux := http.NewServeMux()

	oauth := &oauth2.Config{
		RedirectURL:  cfg.GetHostURL() + "/callback",
		ClientID:     cfg.GithubClientID,
		ClientSecret: cfg.GithubClientSecret,
		Endpoint:     github.Endpoint,
	}

	addRoutes(mux, session, oauth, database)

	handler := middleware.LoggingMiddleware()(mux)
	handler = session.LoadAndSave(handler)

	server := &http.Server{
		Addr:              net.JoinHostPort(cfg.HostIP, cfg.Port),
		Handler:           handler,
		ErrorLog:          slog.NewLogLogger(slog.Default().Handler(), cfg.LogLevel),
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
