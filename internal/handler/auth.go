package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/alexedwards/scs/v2"
	"github.com/confy/godo/internal/db"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func HandleAuthLogin(oauth *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redirect to the OAuth2 login page
		state := uuid.NewString()
		http.Redirect(w, r, oauth.AuthCodeURL(state), http.StatusSeeOther)

	}
}

func HandleAuthLogout(session *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear the session and redirect to the home page
		session.Destroy(r.Context())
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HandleAuthCallback(
	session *scs.SessionManager,
	oauth *oauth2.Config,
	database *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the code from the query string
		code := r.URL.Query().Get("code")
		// Exchange the code for a token
		token, err := oauth.Exchange(context.Background(), code)
		if err != nil {
			errorHandler(w, r, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		// Get the user information
		user, err := getUserInfo(oauth.Client(context.Background(), token))
		if err != nil {
			errorHandler(w, r, "Failed to get user information", http.StatusInternalServerError)
			return
		}

		// Create or get the user in the database
		dbUser, err := db.CreateOrGetUser(context.Background(), database, user)
		if err != nil {
			errorHandler(w, r, "Failed to create or get user", http.StatusInternalServerError)
			return
		}

		// Store user ID and token in the session
		session.Put(r.Context(), "userID", dbUser.ID)
		session.Put(r.Context(), "token", token.AccessToken)

		// Redirect to the original URL or home page
		redirectURL := "/"
		originalURL := session.PopString(r.Context(), "redirect")
		if originalURL != "" {
			redirectURL = originalURL
		}
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func getUserInfo(client *http.Client) (db.CreateUserParams, error) {
	// resp, err := client.Get("https://api.github.com/user")
	resp, err := client.Do(&http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "api.github.com", Path: "/user"},
	})
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
