package route

import (
	"github.com/gorilla/mux"
	controller "github.com/mystardustcaptain/mattodo/pkg/controller"
)

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", controller.Index).Methods("GET")

	router.HandleFunc("/todo", controller.GetTodos).Methods("GET")
	router.HandleFunc("/todo/{id}", controller.GetTodoById).Methods("GET")
	router.HandleFunc("/todo", controller.CreateTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", controller.UpdateTodoById).Methods("PUT")
	router.HandleFunc("/todo/{id}", controller.DeleteTodoById).Methods("DELETE")

	return router
}
