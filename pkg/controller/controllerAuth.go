package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/auth"
	"github.com/mystardustcaptain/mattodo/pkg/model"
)

// RegisterAuthRoutes registers the authentication routes to the router
// URL: /auth/login?provider=google
// URL: /auth/callback?provider=google
// Other providers: facebook, github
func (c *Controller) RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth", c.AuthIndex).Methods("GET")
	router.HandleFunc("/auth/login", c.HandleLogin).Methods("GET")
	router.HandleFunc("/auth/callback", c.HandleCallback).Methods("GET")
}

// HandleLogin initiates the OAuth login process for a given provider
// It redirects the user to the provider's login page
func (c *Controller) HandleLogin(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	config, ok := auth.OAuthConfigs[provider]
	if !ok {
		log.Printf("Unknown OAuth provider")
		respondWithError(w, http.StatusBadRequest, "Unknown OAuth provider")
		return
	}

	url := config.AuthCodeURL(auth.OAuthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleCallback handles the callback from the OAuth provider
// It exchanges the OAuth code for an access token
// and then exchanges the access token for user info.
// User info is made sure available in the database.
// If not, create a new user entry in the database.
// Finally, it creates a JWT token and returns it to the user.
func (c *Controller) HandleCallback(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	state := r.FormValue("state")
	if state != auth.OAuthStateString {
		log.Printf("Invalid state parameter")
		respondWithError(w, http.StatusBadRequest, "Invalid state parameter")
		return
	}

	// get userinfo from token exchanged from OAuth code
	code := r.FormValue("code")
	userInfo, err := auth.GetUserFromOAuthCode(provider, code)
	if err != nil {
		log.Printf("Failed to get user info: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// check if userInfo.Email is a valid email address format
	if !auth.IsEmailValid(userInfo.Email) {
		log.Printf("Invalid email address found from OAuth provider: %s", userInfo.Email)
		respondWithError(w, http.StatusBadRequest, "Invalid email address found from OAuth provider")
		return
	}

	uc := model.UserCollection{DB: c.Database}

	// Try getting user from the database
	user, _ := uc.GetUserByEmail(userInfo.Email)
	if err != nil {
		log.Printf("Failed to get user entry: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if user == nil {
		// If not exist, create a user entry in the database
		log.Printf("User not found, registering user entry.")

		user = &model.User{OAuthProvider: provider, OAuthID: userInfo.ID, Name: userInfo.Name, Email: userInfo.Email}
		if err := uc.CreateUser(user); err != nil {
			log.Printf("Failed to create user entry: %s", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		log.Printf("User entry created for %s", userInfo.Email)
	}

	// TODO: Choosing EMAIL as the checking condition is a simple approach.
	// Unhandled Scenario like, user authenticate using one Provider previously (OAuthProvider + OAuthID)
	// But changed the associated EMAIL, ...

	// Create a JWT token valid for 1 hour
	token, err := auth.CreateToken(user.Email, user.ID, 1)
	if err != nil {
		log.Printf("Failed to create token: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	// Return the token to the user
	respondWithJSON(w, http.StatusOK, token)
}

func (c *Controller) AuthIndex(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Auth Index"})
}
