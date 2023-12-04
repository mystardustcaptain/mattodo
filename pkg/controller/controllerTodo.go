package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/auth"
	"github.com/mystardustcaptain/mattodo/pkg/model"
)

func (c *Controller) GetTodos(w http.ResponseWriter, r *http.Request) {
	// Retrieve iam from the request context
	iam, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("userID: ", iam)
		fmt.Println("Failed to read context")
		respondWithError(w, http.StatusInternalServerError, "Failed to read context")
		return
	}

	todoItems, err := model.GetAllTodoItems(c.Database, iam)
	if err != nil {
		fmt.Println("Failed to get todo items: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, todoItems)
}

func (c *Controller) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// Retrieve iam from the request context
	iam, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("userID: ", iam)
		fmt.Println("Failed to read context")
		respondWithError(w, http.StatusInternalServerError, "Failed to read context")
		return
	}

	var t model.TodoItem
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &t)

	err := t.CreateTodoItem(c.Database, iam)
	if err != nil {
		fmt.Println("Failed to create todo item: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to create todo item: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (c *Controller) DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	// Retrieve iam from the request context
	iam, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("Failed to read context")
		respondWithError(w, http.StatusInternalServerError, "Failed to read context")
		return
	}

	t := model.TodoItem{UserID: iam}

	// Retrieve the todo item id from the request path
	vars := mux.Vars(r)
	todoItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}
	t.ID = todoItemID

	err = t.DeleteTodoItem(c.Database, iam)
	if err != nil {
		fmt.Println("Failed to delete todo item: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to delete todo item: "+err.Error())
		return
	}

	// Operation was successful, but no content to return
	respondWithJSON(w, http.StatusNoContent, nil)
}

func (c *Controller) MarkTodoCompleteById(w http.ResponseWriter, r *http.Request) {
	// Retrieve iam from the request context
	iam, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("Failed to read context")
		respondWithError(w, http.StatusInternalServerError, "Failed to read context")
		return
	}

	t := model.TodoItem{UserID: iam}

	// Retrieve the todo item id from the request path
	vars := mux.Vars(r)
	todoItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}
	t.ID = todoItemID

	err = t.MarkComplete(c.Database, iam)
	if err != nil {
		fmt.Println("Failed to mark complete todo item: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to mark complete todo item: "+err.Error())
		return
	}

	err = t.GetTodoItem(c.Database, iam)
	if err != nil {
		fmt.Println("Failed to retrieve item after mark complete: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve item after mark complete: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (c *Controller) RegisterTodoRoutes(router *mux.Router) {
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(c.GetTodos))).Methods("GET")
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(c.CreateTodo))).Methods("POST")
	router.Handle("/todo/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(c.DeleteTodoById))).Methods("DELETE")
	router.Handle("/todo/{id}/complete", auth.ValidateTokenMiddleware(http.HandlerFunc(c.MarkTodoCompleteById))).Methods("PUT")

}
