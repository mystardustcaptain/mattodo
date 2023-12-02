package controller

import (
	"fmt"
	"net/http"
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
