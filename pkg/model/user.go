package model

import (
	"database/sql"
	"log"
)

// User represents a user of the TODO application, identified by an external OAuth provider.
type User struct {
	ID            int    `json:"id"`
	OAuthProvider string `json:"oauth_provider"`
	OAuthID       string `json:"oauth_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
}

type UserCollection struct {
	DB *sql.DB
}

func (uc *UserCollection) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, oauth_provider, oauth_id, name, email FROM users WHERE email = ?"

	u := User{}
	err := uc.DB.QueryRow(query, email).Scan(&u.ID, &u.OAuthProvider, &u.OAuthID, &u.Name, &u.Email)

	if err != nil {
		log.Printf("Failed to get user by email: %s", err.Error())
		return nil, err
	}

	return &u, nil
}

// expect u to be modified
func (uc *UserCollection) CreateUser(u *User) error {
	query := "INSERT INTO users (oauth_provider, oauth_id, name, email) VALUES (?, ?, ?, ?)"
	res, err := uc.DB.Exec(query, u.OAuthProvider, u.OAuthID, u.Name, u.Email)
	if err != nil {
		log.Printf("Failed to create user: %s", err.Error())
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last insert id: %s", err.Error())
		return err
	}
	u.ID = int(id)

	return err
}

// TODO: There is currently no mechanism to DeleteUser and UpdateUser
