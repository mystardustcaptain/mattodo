package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	_ "github.com/mystardustcaptain/mattodo/pkg/config"
)

// UserInfo represents the user's information returned from the OAuth provider
type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
}

// OAuthConfigurations for multiple providers
var OAuthConfigs map[string]*oauth2.Config

// Initialize OAuth configurations for multiple providers
// required environment variables
func init() {
	OAuthConfigs = map[string]*oauth2.Config{
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
			Scopes:       []string{"email, name"},
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
}

// OAuthStateString is a randomly generated string to protect against CSRF attacks
var OAuthStateString = "wzp-bdt*czm8GEQ9kuc"

// GetUserFromOAuthCode exchanges an OAuth code for a token, then fetches user information
// provider: google, facebook, github
// code: auth code returned from the OAuth provider
// returns the user information or an error
func GetUserFromOAuthCode(provider string, code string) (*UserInfo, error) {
	config, ok := OAuthConfigs[provider]
	if !ok {
		log.Printf("Unknown OAuth provider: %s\n", provider)
		return nil, errors.New("unknown OAuth provider")
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %s\n", err.Error())
		return nil, err
	}

	userInfo, err := fetchUserInfo(provider, token.AccessToken)
	if err != nil {
		log.Printf("Failed to fetch user info: %s\n", err.Error())
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
		log.Printf("Failed to create request: %s\n", err.Error())
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
		log.Printf("Failed to make request: %s\n", err.Error())
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body: %s\n", err.Error())
		return nil, err
	}

	// Unmarshal the JSON data into the UserInfo struct
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		log.Printf("Failed to unmarshal JSON: %s\n", err.Error())
		return nil, err
	}

	return &userInfo, nil
}

// CreateToken creates a JWT token with the userEmail as the subject
// returns the token or an error
// param userEmail: the user's email address
// param userID: the user's ID in database
// param hour: the validity of the token in hour
func CreateToken(userEmail string, userID int, hour int) (string, error) {
	var mySigningKey = []byte(os.Getenv("SIGNING_KEY"))

	// use HASH256 to sign the token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	// token valid for x hour
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(hour)).Unix()
	// info about the user to be encoded in the token
	claims["userEmail"] = userEmail
	claims["userID"] = userID

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Printf("Failed to sign token: %s\n", err.Error())
		return "", err
	}

	return tokenString, nil
}

// ValidateTokenMiddleware validates the token from the Authorization header
// every request with this middleware will require a valid token
// Note: only appllies to routes that require authentication
func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)

		if tokenString == "" {
			log.Printf("Authorization token is required\n")
			http.Error(w, "Authorization token is required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v\n", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SIGNING_KEY")), nil
		})

		if err != nil {
			log.Printf("Failed to parse token: %s\n", err.Error())
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// The token is valid and not expired
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract the user info from the token
			userEmail := claims["userEmail"].(string)
			userID := claims["userID"].(float64)

			if userEmail == "" || userID <= 0 {
				// Handle error: userEmail or userID not found in token
				log.Printf("userEmail or userID not found in the token\n")
				http.Error(w, "userEmail or userID not found in the token", http.StatusUnauthorized)
				return
			}

			// Add the db userID to the request context
			ctx := context.WithValue(r.Context(), "userID", int(userID))
			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			log.Printf("Invalid authorization token\n")
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
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
