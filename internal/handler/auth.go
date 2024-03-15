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
		// Redirect to the OAuth2 login page
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

		// Get the user information
		user, err := getUserInfo(oauthConfig.Client(context.Background(), token))
		if err != nil {
			http.Error(w, "Failed to get user information", http.StatusInternalServerError)
			return
		}

		// Create or get the user in the database
		dbUser, err := db.CreateOrGetUser(context.Background(), dbQueries, user)
		if err != nil {
			http.Error(w, "Failed to create or get user", http.StatusInternalServerError)
			return
		}

		// Store user ID and token in the session
		sessionManager.Put(r.Context(), "user_id", dbUser.ID)
		sessionManager.Put(r.Context(), "token", token.AccessToken)

		// Redirect to the original URL or home page
		redirectURL := "/"
		originalURL := sessionManager.PopString(r.Context(), "redirect")
		if originalURL != "" {
			redirectURL = originalURL
		}
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func getUserInfo(client *http.Client) (db.CreateUserParams, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return db.CreateUserParams{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return db.CreateUserParams{}, err
	}

	var user db.CreateUserParams
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return db.CreateUserParams{}, err
	}

	return user, nil
}
