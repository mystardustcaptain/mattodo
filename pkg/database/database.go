package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database
// dbType: sqlite3, mysql, postgres
// dbPath: path to the database file
func InitDB(dbType string, dbPath string) *sql.DB {
	// Initialize database
	db, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
