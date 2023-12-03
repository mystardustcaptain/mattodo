package controller

import (
	"fmt"
	"net/http"

	"github.com/mystardustcaptain/mattodo/pkg/auth"

	"github.com/gorilla/mux"
)

// HandleLogin initiates the OAuth login process for a given provider
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	config, ok := auth.OAuthConfigs[provider]
	if !ok {
		http.Error(w, "Unknown OAuth provider", http.StatusBadRequest)
		return
	}

	url := config.AuthCodeURL(auth.OAuthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleCallback handles the callback from the OAuth provider
func HandleCallback(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	state := r.FormValue("state")
	if state != auth.OAuthStateString {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	userInfo, err := auth.GetUserFromOAuthToken(provider, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Todo:
	// create a session or JWT token for the authenticated user??
	// create or update the user record in database??

	fmt.Fprintf(w, "User Info: %+v", userInfo)
}

func AuthIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth Index")
}

// RegisterAuthRoutes registers the authentication routes to the router
// URL: /auth/login?provider=google
// URL: /auth/callback?provider=google
// Other providers: facebook, github
func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth", AuthIndex).Methods("GET")
	router.HandleFunc("/auth/login", HandleLogin).Methods("GET")
	router.HandleFunc("/auth/callback", HandleCallback).Methods("GET")
}
