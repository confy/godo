package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func RequireLogin(next http.Handler, session *scs.SessionManager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := session.GetInt64(r.Context(), "userID")
		// Check if user is authenticated
		if userID == 0 {
			// Save the current URL in the session so we can redirect after login
			session.Put(r.Context(), "redirect", r.URL.Path)
			// Redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// User is authenticated, call the next handler
		next.ServeHTTP(w, r)
	})
}
