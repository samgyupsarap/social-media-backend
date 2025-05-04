package routes

import (
	"backend/controllers"
	"backend/db/socmed"
	"backend/handlers"
	"backend/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(queries *socmed.Queries) *mux.Router {
	r := mux.NewRouter()

	loginController := controllers.NewLoginController(queries)
	userHandler := handlers.NewUserHandlerController(queries)

	r.HandleFunc("/api/login", loginController.Login).Methods("POST")
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.CheckToken())

	
	api.HandleFunc("/user", userHandler.UserHandler)

	return r
}
