package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/auth"
	"github.com/mystardustcaptain/mattodo/pkg/model"
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

	// Check if user email already exists in the database
	exist, err := model.IsUserExistByEmail(c.Database, userInfo.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If not exist, create a user entry in the database
	if !exist {
		fmt.Println("User not found, registering user entry.")

		db_user := model.NewUser(provider, userInfo.ID, userInfo.Name, userInfo.Email)
		if err := db_user.CreateUser(c.Database); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("User entry created for ", userInfo.Email)
	}

	// TODO: Choosing EMAIL as the checking condition is a simple approach.
	// Unhandled Scenario like, user authenticate using one Provider previously (OAuthProvider + OAuthID)
	// But changed the associated EMAIL, ...

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
