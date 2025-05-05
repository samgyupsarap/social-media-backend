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
	commentHandler := handlers.NewCommentHandlerController(queries)
	postHandler := handlers.NewPostHandlerController(queries)

	r.HandleFunc("/api/login", loginController.Login).Methods("POST")
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.CheckToken())

	api.HandleFunc("/user", userHandler.UserHandler)
	api.HandleFunc("/comment", commentHandler.CommentHandler)
	api.HandleFunc("/post", postHandler.PostHandler)

	return r
}
