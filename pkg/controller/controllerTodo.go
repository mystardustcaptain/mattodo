package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/auth"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from the request context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		fmt.Println("Failed to read context")
		
	}
	fmt.Fprintf(w, "Get Todos for user %s", userID)
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get Todo By Id")
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create Todo")
}

func UpdateTodoById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Update Todo By Id")
}

func DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete Todo By Id")
}

func RegisterTodoRoutes(router *mux.Router) {
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(GetTodos))).Methods("GET")
	router.Handle("/todo/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(GetTodoById))).Methods("GET")
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(CreateTodo))).Methods("POST")
	router.Handle("/todo/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(UpdateTodoById))).Methods("PUT")
	router.Handle("/todo/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(DeleteTodoById))).Methods("DELETE")
}
