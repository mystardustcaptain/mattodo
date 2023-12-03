package route

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/controller"
)

func InitializeRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	c := controller.NewController(db)

	c.RegisterRoutes(router)
	c.RegisterTodoRoutes(router)
	c.RegisterAuthRoutes(router)

	return router
}
