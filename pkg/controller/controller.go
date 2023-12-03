package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Controller struct {
	Database *sql.DB
}

func NewController(db *sql.DB) *Controller {
	return &Controller{
		Database: db,
	}
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up and running")
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", c.Index).Methods("GET")
}
