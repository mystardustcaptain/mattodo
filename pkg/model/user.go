package model

// User represents a user of the TODO application, identified by an external OAuth provider.
type User struct {
	ID            int    `json:"id"`
	OAuthProvider string `json:"oauthProvider"`
	OAuthID       string `json:"oauthId"`
	Name          string `json:"name"`
	Email         string `json:"email"`
}

// NewUser creates a new User instance with details from an OAuth provider.
func NewUser(oAuthProvider, oAuthID, name, email string) *User {
	return &User{
		OAuthProvider: oAuthProvider,
		OAuthID:       oAuthID,
		Name:          name,
		Email:         email,
	}
}
