package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// OAuthConfigurations for multiple providers
var OAuthConfigs = map[string]*oauth2.Config{
	"google": {
		RedirectURL:  "http://localhost:9003/auth/callback?provider=google",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	},
	"facebook": {
		RedirectURL:  "http://localhost:9003/auth/callback?provider=facebook",
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
// provider: google, facebook, github
// code: code returned from the OAuth provider
// returns the user information or an error
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
// provider: google, facebook, github
// accessToken: token returned from the OAuth provider
// returns the user information or an error
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

// CreateToken creates a JWT token with the userEmail as the subject
func CreateToken(userEmail string) (string, error) {
	var mySigningKey = []byte(os.Getenv("SIGNING_KEY"))

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["userEmail"] = userEmail
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)

		if tokenString == "" {
			http.Error(w, "Authorization token is required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SIGNING_KEY")), nil
		})

		if err != nil {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userEmail, ok := claims["userEmail"].(string); ok && userEmail != "" {
				// Add the userEmail to the request context
				ctx := context.WithValue(r.Context(), "userEmail", userEmail)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				// Handle error: userEmail not found in token
				http.Error(w, "userEmail not found in the token", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
		}
	})
}

// extractToken extracts the token from the Authorization header
// expected format:
// Authorization: Bearer {token-body}
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

// UserInfo represents the user's information returned from the OAuth provider
type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
}
