package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func requireLogin(next http.Handler, sessionManager *scs.SessionManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_id := sessionManager.GetString(r.Context(), "user_id")
		// Check if user is authenticated
		if user_id == "" {
			// Redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// User is authenticated, call the next handler
		next.ServeHTTP(w, r)
	})
}
