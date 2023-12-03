package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/auth"
)

// HandleLogin initiates the OAuth login process for a given provider
func (c *Controller) HandleLogin(w http.ResponseWriter, r *http.Request) {
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
func (c *Controller) HandleCallback(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	state := r.FormValue("state")
	if state != auth.OAuthStateString {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// get userinfo from token exchanged from OAuth code
	code := r.FormValue("code")
	userInfo, err := auth.GetUserFromOAuthToken(provider, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Todo:
	// create or update the user record in database??

	// Create a JWT token valid for 1 hour
	token, err := auth.CreateToken(userInfo.Email)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Return the token to the user
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Token: %s\n", token)
	fmt.Fprintf(w, "User Info: %+v\n", userInfo)
}

func (c *Controller) AuthIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth Index")
}

// RegisterAuthRoutes registers the authentication routes to the router
// URL: /auth/login?provider=google
// URL: /auth/callback?provider=google
// Other providers: facebook, github
func (c *Controller) RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth", c.AuthIndex).Methods("GET")
	router.HandleFunc("/auth/login", c.HandleLogin).Methods("GET")
	router.HandleFunc("/auth/callback", c.HandleCallback).Methods("GET")
}
