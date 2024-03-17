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

		todo, err := database.CreateTodo(context.Background(), db.CreateTodoParams{
			UserID: user.ID,
			Title:   "Test todo",
			Description: sql.NullString{String: "This is a test todo", Valid: true},
		})

		templ.Handler(views.TestPage(user, todo)).ServeHTTP(w, r)
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templ.Handler(views.IndexPage("Hello, world!")).ServeHTTP(w, r)

}
