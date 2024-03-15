package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/confy/godo/internal/db"
	"golang.org/x/oauth2"
)

func HandleAuthLogin(oauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redirect to the oauth2 login page
		http.Redirect(w, r, oauthConfig.AuthCodeURL("state"), http.StatusSeeOther)
	}
}

func HandleAuthCallback(sessionManager *scs.SessionManager, oauthConfig *oauth2.Config, dbQueries *db.Queries) http.HandlerFunc {
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
		dbUser, err := db.CreateOrGetUser(context.Background(), dbQueries, user)

		if err != nil {
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}

		sessionManager.Put(r.Context(), "user_id", dbUser.ID)
		sessionManager.Put(r.Context(), "token", token.AccessToken)

		redirectURL := "/"
		originalURL := sessionManager.PopString(r.Context(), "redirect")
		if originalURL != "" {
			redirectURL = originalURL
		}

		// Redirect to the home page
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
