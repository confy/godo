package handler

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/confy/godo/internal/db"
	"github.com/confy/godo/internal/models"
	"github.com/confy/godo/views"
)

func HandleTestPage(database *db.Queries, session *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := session.GetInt64(r.Context(), "userID")
		user, err := database.GetUserById(context.Background(), userID)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}

		// TEST code
		err = testCreateTodo(database, user, w)
		if err != nil {
			return
		}

		todos, err := database.GetTodosByUserId(context.Background(), user.ID)
		if err != nil {
			http.Error(w, "Failed to get todos", http.StatusInternalServerError)
			return
		}

		displayUser := models.DisplayUser{
			LoggedIn:  true,
			Login:     user.Login,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		}
		templ.Handler(views.TestPage(displayUser, todos)).ServeHTTP(w, r)
	}
}

func testCreateTodo(database *db.Queries, user db.User, w http.ResponseWriter) error {
	_, err := database.CreateTodo(context.Background(), db.CreateTodoParams{
		UserID:      user.ID,
		Title:       "Test todo",
		Description: sql.NullString{String: "This is a test todo", Valid: true},
	})
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return err
	}

	_, err = database.CreateTodo(context.Background(), db.CreateTodoParams{
		UserID:      user.ID,
		Title:       "Test todo",
		Description: sql.NullString{String: "This is a test todo", Valid: true},
	})
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return err
	}
	return err
}

func HandleRoot(database *db.Queries, session *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			if _, err := os.Stat(filepath.Join("static", r.URL.Path)); err == nil {
				http.ServeFile(w, r, "static"+r.URL.Path)
				return
			}
			errorHandler(w, r, http.StatusNotFound)
			return

		}

		// Root is the only page that doesn't require a user to be logged in
		var displayUser models.DisplayUser
		userID := session.GetInt64(r.Context(), "userID")
		if userID == 0 {
			displayUser = models.DisplayUser{
				LoggedIn: false,
			}
		} else {
			user, err := database.GetUserById(context.Background(), userID)
			if err != nil {
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			}
			displayUser = models.DisplayUser{
				LoggedIn:  true,
				Login:     user.Login,
				Email:     user.Email,
				AvatarURL: user.AvatarURL,
			}
		}
		templ.Handler(views.IndexPage("Hello, world!", displayUser)).ServeHTTP(w, r)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	w.Write([]byte("fancy custom 404 page!"))
}
