package main

import (
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/mystardustcaptain/mattodo/pkg/config"
	"github.com/mystardustcaptain/mattodo/pkg/database"
	"github.com/mystardustcaptain/mattodo/pkg/route"
)

func main() {
	// Read configuration
	port := os.Getenv("SERVICE_PORT")

	// Initialize database
	db := database.InitDB(os.Getenv("DB_TYPE"), os.Getenv("DB_PATH"))

	// Initialize router
	r := route.InitializeRoutes(db)

	// Start web server
	http.ListenAndServe(port, r)
}
