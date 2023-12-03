package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up and running")
}

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", Index).Methods("GET")
}
