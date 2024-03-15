package handler

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/confy/godo/internal/db"
	"github.com/confy/godo/views"
)

func HandleTestPage(dbQueries *db.Queries, sessionManager *scs.SessionManager) http.HandlerFunc {
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

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templ.Handler(views.IndexPage("Hello, world!")).ServeHTTP(w, r)

}
