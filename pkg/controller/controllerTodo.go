package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mystardustcaptain/mattodo/pkg/auth"
	"github.com/mystardustcaptain/mattodo/pkg/model"
)

// Register routes for the controller related to todo items
func (c *Controller) RegisterTodoRoutes(router *mux.Router) {
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(c.GetTodos))).Methods("GET")
	router.Handle("/todo", auth.ValidateTokenMiddleware(http.HandlerFunc(c.CreateTodo))).Methods("POST")
	router.Handle("/todo/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(c.DeleteTodoById))).Methods("DELETE")
	router.Handle("/todo/{id}/complete", auth.ValidateTokenMiddleware(http.HandlerFunc(c.MarkTodoCompleteById))).Methods("PUT")
}

// GetTodos retrieves all todo items for the authenticated user
// with userID saved in the request context
func (c *Controller) GetTodos(w http.ResponseWriter, r *http.Request) {
	// Retrieve iam / db userID from the request context
	iam, ok := r.Context().Value(auth.ContextUserIDKey).(int)
	if !ok {
		log.Printf("Failed to read context")
		respondWithError(w, http.StatusInternalServerError, "Failed to read context")
		return
	}

	// Retrieve all todo items for the user
	todoItems, err := model.GetAllTodoItems(c.Database, iam)
	if err != nil {
		log.Printf("Failed to get all todo items: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, todoItems)
}

// CreateTodo creates a new todo item for the authenticated user
// with userID saved in the request context
func (c *Controller) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// Retrieve iam from the request context
	iam, ok := r.Context().Value(auth.ContextUserIDKey).(int)
	if !ok {
		log.Printf("Failed to read context")
		respondWithError(w, http.StatusInternalServerError, "Failed to read context")
		return
	}

	var t model.TodoItem

	// Decode the request body into a TodoItem struct
	// content provided will be used to create a new todo item
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &t)

	// Create the todo item in the database
	err := t.CreateTodoItem(c.Database, iam)
	if err != nil {
		log.Printf("Failed to create todo item: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to create todo item: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

// DeleteTodoById deletes a todo item for the authenticated user
// with userID saved in the request context
func (c *Controller) DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	// Retrieve iam from the request context
	iam, ok := r.Context().Value(auth.ContextUserIDKey).(int)
	if !ok {
		log.Printf("Failed to read context")
		respondWithError(w, http.StatusInternalServerError, "Failed to read context")
		return
	}

	t := model.TodoItem{UserID: iam}

	// Retrieve the todo item id from the request path
	// This is the target todo item to be deleted
	vars := mux.Vars(r)
	todoItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}
	t.ID = todoItemID

	// Delete the todo item from the database
	err = t.DeleteTodoItem(c.Database, iam)
	if err != nil {
		log.Printf("Failed to delete todo item: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to delete todo item: "+err.Error())
		return
	}

	// Operation was successful, but no content to return
	respondWithJSON(w, http.StatusNoContent, nil)
}

// MarkTodoCompleteById marks a todo item as complete for the authenticated user
// with userID saved in the request context
func (c *Controller) MarkTodoCompleteById(w http.ResponseWriter, r *http.Request) {
	// Retrieve iam from the request context
	iam, ok := r.Context().Value(auth.ContextUserIDKey).(int)
	if !ok {
		log.Printf("Failed to read context")
		respondWithError(w, http.StatusInternalServerError, "Failed to read context")
		return
	}

	t := model.TodoItem{UserID: iam}

	// Retrieve the todo item id from the request path
	// This is the target todo item to be marked complete
	vars := mux.Vars(r)
	todoItemID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Invalid todo ID")
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}
	t.ID = todoItemID

	// Mark the todo item as complete in the database
	err = t.MarkComplete(c.Database, iam)
	if err != nil {
		log.Printf("Failed to mark complete todo item: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to mark complete todo item: "+err.Error())
		return
	}

	// Retrieve the todo item from the database
	// to return to the user
	err = t.GetTodoItem(c.Database, iam)
	if err != nil {
		log.Printf("Failed to retrieve item after mark complete: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve item after mark complete: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}
