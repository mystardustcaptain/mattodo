package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// InitDB initializes the database
// dbType: sqlite, mysql, postgres
// dbPath: path to the database file
func InitDB(dbType string, dbPath string) *sql.DB {
	// Initialize database
	db, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
