package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	router "github.com/mystardustcaptain/mattodo/pkg/route"
)

func main() {
	// Load environment variables from .env file
	// if running locally
	if os.Getenv("DOCKER_ENV_SET") != "true" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

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
	r := router.InitializeRoutes(db)

	// Start web server
	http.ListenAndServe(port, r)
}
