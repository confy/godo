package server

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// CreateSessionManager creates a new session manager with some default settings and optional https
func CreateSessionManager(secure bool) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Secure = secure
	sessionManager.Cookie.HttpOnly = true

	return sessionManager
}
