package handler

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
			errorHandler(w, r, "Failed to get todos", http.StatusInternalServerError)
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
			errorHandler(w, r, "Page not found", http.StatusNotFound)
			return
		}
		// Root is the only page that doesn't require a user to be logged in, but we still want to display the user if they are logged in
		var displayUser models.DisplayUser
		userID := session.GetInt64(r.Context(), "userID")
		displayUser, err := newFunction(userID, displayUser, database, w)
		if err != nil {
			errorHandler(w, r, "Unable to get user from db",http.StatusInternalServerError)
		}
		templ.Handler(views.IndexPage("Hello, world!", displayUser)).ServeHTTP(w, r)
	}
}

func newFunction(userID int64, displayUser models.DisplayUser, database *db.Queries, w http.ResponseWriter) (models.DisplayUser, error) {
	if userID == 0 {
		displayUser = models.DisplayUser{
			LoggedIn: false,
		}
	} else {
		user, err := database.GetUserById(context.Background(), userID)
		if err != nil {
			return models.DisplayUser{}, err
		}
		displayUser = models.DisplayUser{
			LoggedIn:  true,
			Login:     user.Login,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		}
	}
	return displayUser, nil
}

func errorHandler(w http.ResponseWriter, _ *http.Request, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(strconv.Itoa(status) + " " + message))	
}
