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
	msController := controllers.NewMicrosoftController(queries)

	r.HandleFunc("/api/login", loginController.Login).Methods("POST")
	auth := r.PathPrefix("/api/auth").Subrouter()
	auth.HandleFunc("/microsoft/login", msController.LoginHandler).Methods("GET")
	auth.HandleFunc("/microsoft/callback", msController.CallbackHandler).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.CheckToken())

	api.HandleFunc("/user", userHandler.UserHandler)
	api.HandleFunc("/comment", commentHandler.CommentHandler)
	api.HandleFunc("/post", postHandler.PostHandler)

	return r
}
