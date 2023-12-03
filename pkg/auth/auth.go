package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// OAuthConfigurations for multiple providers
var OAuthConfigs = map[string]*oauth2.Config{
	"google": {
		RedirectURL:  "http://localhost:8080//auth/callback?provider=google",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	},
	"facebook": {
		RedirectURL:  "http://localhost:8080/auth/callback?provider=facebook",
		ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		Scopes:       []string{"email"},
		Endpoint:     facebook.Endpoint,
	},
	"github": {
		RedirectURL:  "http://localhost:9003/auth/callback?provider=github",
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	},
}

// OAuthStateString is a randomly generated string to protect against CSRF attacks
// Replace with a random or dynamically generated string
var OAuthStateString = "random"

// GetUserFromOAuthToken exchanges an OAuth code for a token, then fetches user information
func GetUserFromOAuthToken(provider string, code string) (*UserInfo, error) {
	config, ok := OAuthConfigs[provider]
	if !ok {
		return nil, errors.New("unknown OAuth provider")
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	userInfo, err := fetchUserInfo(provider, token.AccessToken)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// fetchUserInfo retrieves the user information from the OAuth provider
func fetchUserInfo(provider, accessToken string) (*UserInfo, error) {
	var endpoint string
	switch provider {
	case "google":
		endpoint = "https://www.googleapis.com/oauth2/v2/userinfo"
	case "facebook":
		endpoint = "https://graph.facebook.com/me?fields=id,name,email"
	case "github":
		endpoint = "https://api.github.com/user"
	default:
		return nil, errors.New("unknown OAuth provider for user info")
	}

	// Create a new request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header for GitHub
	if provider == "github" {
		req.Header.Set("Authorization", "token "+accessToken)
	} else {
		// For others, use access token in query params
		q := req.URL.Query()
		q.Add("access_token", accessToken)
		req.URL.RawQuery = q.Encode()
	}

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the UserInfo struct
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// UserInfo represents the user's information returned from the OAuth provider
type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
}
