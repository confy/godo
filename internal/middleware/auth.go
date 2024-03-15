package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func RequireLogin(next http.Handler, sessionManager *scs.SessionManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_id := sessionManager.GetInt64(r.Context(), "user_id")
		// Check if user is authenticated
		if user_id == 0 {
			// Save the current URL in the session so we can redirect after login
			sessionManager.Put(r.Context(), "redirect", r.URL.Path)
			// Redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// User is authenticated, call the next handler
		next.ServeHTTP(w, r)
	})
}
