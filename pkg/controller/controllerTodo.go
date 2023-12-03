package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get Todos")
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
	router.HandleFunc("/todo", GetTodos).Methods("GET")
	router.HandleFunc("/todo/{id}", GetTodoById).Methods("GET")
	router.HandleFunc("/todo", CreateTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", UpdateTodoById).Methods("PUT")
	router.HandleFunc("/todo/{id}", DeleteTodoById).Methods("DELETE")

}
