package controller

import (
	"database/sql"
	"encoding/json"
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

// RegisterRoutes registers routes for the controller
func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", c.Index).Methods("GET")
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Server is up and running"})
}

// respondWithError is a helper function to respond with an error and a status code
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON is a helper function to respond with JSON and a status code
// payload can be nil
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	if payload != nil {
		response, _ := json.Marshal(payload)
		// sequence of these lines matters
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	} else {
		w.WriteHeader(code)
	}
}
