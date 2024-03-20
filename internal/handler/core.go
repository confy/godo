package handler

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/confy/godo/internal/db"
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

		templ.Handler(views.TestPage(user, todos)).ServeHTTP(w, r)
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

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templ.Handler(views.IndexPage("Hello, world!")).ServeHTTP(w, r)

}
