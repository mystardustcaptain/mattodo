package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	// db "github.com/mystardustcaptain/mattodo/pkg/db"
	router "github.com/mystardustcaptain/mattodo/pkg/route"
)

// App struct with router and database instances plus other configurations
type App struct {
	Router   *mux.Router
	Database *sql.DB
	Port     string
	DBType   string
	DBPath   string
}

func main() {
	a := App{}

	// Load environment variables from .env file
	// if running locally
	if os.Getenv("DOCKER_ENV_SET") != "true" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Read configuration
	a.Port = os.Getenv("SERVICE_PORT")
	a.DBType = os.Getenv("DB_TYPE")
	a.DBPath = os.Getenv("DB_PATH")

	// Initialize database
	// a.Database = db.InitDB(a.DBType, a.DBPath)

	// Initialize router
	a.Router = router.InitializeRoutes()

	// Start web server
	http.ListenAndServe(a.Port, a.Router)
}
