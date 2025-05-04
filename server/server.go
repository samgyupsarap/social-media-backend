package server

import (
	"backend/db/socmed"
	"backend/routes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func StartServer(db *sql.DB) {
	queries := socmed.New(db)

	router := routes.SetupRoutes(queries)

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	}).Methods(http.MethodGet)

	port := ":8080"
	fmt.Println("Starting server at", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
