package route

import (
	"github.com/gorilla/mux"
	controller "github.com/mystardustcaptain/mattodo/pkg/controller"
)

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	controller.RegisterRoutes(router)
	controller.RegisterTodoRoutes(router)
	controller.RegisterAuthRoutes(router)

	return router
}
