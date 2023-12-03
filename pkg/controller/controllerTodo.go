package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/auth"
)

func (c *Controller) GetTodos(w http.ResponseWriter, r *http.Request) {
	// Retrieve userEmail from the request context
	userEmail, ok := r.Context().Value("userEmail").(string)
	if !ok {
		fmt.Println("Failed to read context")

	}
	fmt.Fprintf(w, "Get Todos for user %s", userEmail)
}

func (c *Controller) GetTodoById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get Todo By Id")
}

func (c *Controller) CreateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create Todo")
}

func (c *Controller) UpdateTodoById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Update Todo By Id")
}

func (c *Controller) DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete Todo By Id")
}

func (c *Controller) RegisterTodoRoutes(router *mux.Router) {
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(c.GetTodos))).Methods("GET")
	router.Handle("/todo/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(c.GetTodoById))).Methods("GET")
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(c.CreateTodo))).Methods("POST")
	router.Handle("/todo/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(c.UpdateTodoById))).Methods("PUT")
	router.Handle("/todo/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(c.DeleteTodoById))).Methods("DELETE")
}
