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
		fmt.Println("Unknown OAuth provider")
		respondWithError(w, http.StatusBadRequest, "Unknown OAuth provider")
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
		fmt.Println("Invalid state parameter")
		respondWithError(w, http.StatusBadRequest, "Invalid state parameter")
		return
	}

	// get userinfo from token exchanged from OAuth code
	code := r.FormValue("code")
	userInfo, err := auth.GetUserFromOAuthToken(provider, code)
	if err != nil {
		fmt.Println("Failed to get user info: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if user email already exists in the database
	exist, err := model.IsUserExistByEmail(c.Database, userInfo.Email)
	if err != nil {
		fmt.Println("Failed to check user existance: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Complete db user data needed
	var db_user model.User

	// If exist, get the user entry from the database
	if exist {
		db_user = model.User{Email: userInfo.Email}
		if err := db_user.GetUserByEmail(c.Database); err != nil {
			fmt.Println("Failed to get user entry: ", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println("User entry found: ", userInfo.Email)
	} else {
		// If not exist, create a user entry in the database
		fmt.Println("User not found, registering user entry.")

		db_user = model.User{OAuthProvider: provider, OAuthID: userInfo.ID, Name: userInfo.Name, Email: userInfo.Email}
		if err := db_user.CreateUser(c.Database); err != nil {
			fmt.Println("Failed to create user entry: ", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println("User entry created for ", userInfo.Email)
	}

	// TODO: Choosing EMAIL as the checking condition is a simple approach.
	// Unhandled Scenario like, user authenticate using one Provider previously (OAuthProvider + OAuthID)
	// But changed the associated EMAIL, ...

	// Create a JWT token valid for 1 hour
	token, err := auth.CreateToken(db_user.Email, db_user.ID, 1)
	if err != nil {
		fmt.Println("Failed to create token: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	// Return the token to the user
	respondWithJSON(w, http.StatusOK, token)
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
