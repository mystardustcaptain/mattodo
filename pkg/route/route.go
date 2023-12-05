package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/controller"
)

// InitializeRoutes initializes the routes for the application.
// Any new routes should be registered here.
func InitializeRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	c := controller.NewController(db)

	c.RegisterRoutes(router)
	c.RegisterTodoRoutes(router)
	c.RegisterAuthRoutes(router)

	return router
}
