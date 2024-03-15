package server

import (
	"context"
	"encoding/json"
	"fmt"
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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func New(logger *slog.Logger, config *Config, dbQueries *db.Queries) *http.Server {
	mux := http.NewServeMux()

	sessionManager := scs.New()
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = config.UseHttps

	oauthConfig := &oauth2.Config{
		RedirectURL:  config.GetHostURL() + "/callback",
		ClientID:     config.GithubClientID,
		ClientSecret: config.GithubClientSecret,
		Endpoint:     github.Endpoint,
	}
	handler := sessionManager.LoadAndSave(mux)
	handler = middleware.LoggingMiddleware(logger)(handler)
	addRoutes(mux, sessionManager, oauthConfig, dbQueries)

	server := &http.Server{
		Addr:     net.JoinHostPort(config.Host, config.Port),
		Handler:  handler,
		ErrorLog: slog.NewLogLogger(logger.Handler(), config.LogLevel),
	}
	return server
}

func addRoutes(mux *http.ServeMux, sessionManager *scs.SessionManager, oauthConfig *oauth2.Config, dbQueries *db.Queries) {
	mux.HandleFunc("/login", handleLogin(oauthConfig))
	mux.HandleFunc("/callback", handleCallback(sessionManager, oauthConfig, dbQueries))
	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("/test", middleware.RequireLogin(handleTestPage(dbQueries, sessionManager), sessionManager))
}

func handleTestPage(dbQueries *db.Queries, sessionManager *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id := sessionManager.Get(r.Context(), "user_id").(int64)
		user, err := dbQueries.GetUserById(context.Background(), user_id)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}
		templ.Handler(views.TestPage(user)).ServeHTTP(w, r)
	}
}

func handleLogin(oauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redirect to the oauth2 login page
		http.Redirect(w, r, oauthConfig.AuthCodeURL("state"), http.StatusSeeOther)
	}
}

func handleCallback(sessionManager *scs.SessionManager, oauthConfig *oauth2.Config, dbQueries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the code from the query string
		code := r.URL.Query().Get("code")
		// Exchange the code for a token
		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}
		// Get the user

		resp, err := oauthConfig.Client(context.Background(), token).Get("https://api.github.com/user")
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return

		}
		decoder := json.NewDecoder(resp.Body)
		var user db.CreateUserParams
		err = decoder.Decode(&user)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}
		fmt.Printf("user: %v", user)

		dbUser, err := db.CreateOrGetUser(context.Background(), dbQueries, user)

		fmt.Printf("dbUser: %v", dbUser)

		if err != nil {
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}

		// Save the user to the session
		fmt.Printf("user: %v", dbUser.ID)
		fmt.Printf("token: %v", token.AccessToken)

		sessionManager.Put(r.Context(), "user_id", dbUser.ID)
		sessionManager.Put(r.Context(), "token", token.AccessToken)

		// Redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
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
