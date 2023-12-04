package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/mystardustcaptain/mattodo/pkg/config"
	"github.com/mystardustcaptain/mattodo/pkg/route"
)

func main() {
	// Read configuration
	port := os.Getenv("SERVICE_PORT")
	dbType := os.Getenv("DB_TYPE")
	dbPath := os.Getenv("DB_PATH")

	// Initialize database
	db, err := sql.Open(dbType, dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize router
	r := route.InitializeRoutes(db)

	// Start web server
	http.ListenAndServe(port, r)
}
