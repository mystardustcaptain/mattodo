package model

import "database/sql"

// User represents a user of the TODO application, identified by an external OAuth provider.
type User struct {
	ID            int    `json:"id"`
	OAuthProvider string `json:"oauth_provider"`
	OAuthID       string `json:"oauth_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
}

func IsUserExistByEmail(db *sql.DB, email string) (bool, error) {
	var exists bool

	// The query checks if there is at least one entry with the given email
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"

	// Execute the query
	err := db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (u *User) GetUserByEmail(db *sql.DB) error {
	query := "SELECT id, oauth_provider, oauth_id, name, email FROM users WHERE email = ?"
	return db.QueryRow(query, u.Email).Scan(&u.ID, &u.OAuthProvider, &u.OAuthID, &u.Name, &u.Email)
}

func (u *User) CreateUser(db *sql.DB) error {
	query := "INSERT INTO users (oauth_provider, oauth_id, name, email) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, u.OAuthProvider, u.OAuthID, u.Name, u.Email)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)

	return err
}

// TODO: There is currently no mechanism to DeleteUser and UpdateUser
