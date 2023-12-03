package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/auth"
	"github.com/mystardustcaptain/mattodo/pkg/model"
)

func (c *Controller) GetTodos(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the request context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("userID: ", userID)
		fmt.Println("Failed to read context")
		respondWithError(w, http.StatusInternalServerError, "Failed to read context")
		return
	}

	todoItems, err := model.GetAllTodoItems(c.Database, userID)
	if err != nil {
		fmt.Println("Failed to get todo items: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, todoItems)
}

func (c *Controller) CreateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create Todo")
}

func (c *Controller) DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete Todo By Id")
}

func (c *Controller) MarkTodoCompleteById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mark Todo to Done By Id")
}

func (c *Controller) RegisterTodoRoutes(router *mux.Router) {
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(c.GetTodos))).Methods("GET")
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(c.CreateTodo))).Methods("POST")
	router.Handle("/todo/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(c.DeleteTodoById))).Methods("DELETE")
	router.Handle("/todo/{id}/complete", auth.ValidateTokenMiddleware(http.HandlerFunc(c.MarkTodoCompleteById))).Methods("PUT")

}
