package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbType, dbPath string) *sql.DB {
	// Initialize database
	db, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}
