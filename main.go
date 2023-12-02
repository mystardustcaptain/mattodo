package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
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

	// Read configuration
	a.Port = ":9003"
	a.DBType = "sqlite3"
	a.DBPath = "./mainDB.sqlite3"

	// Initialize database
	// a.Database = db.InitDB(a.DBType, a.DBPath)

	// Initialize router
	a.Router = router.InitializeRoutes()

	// Start web server
	http.ListenAndServe(a.Port, a.Router)
}
