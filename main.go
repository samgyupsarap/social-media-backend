package main

import (
	"backend/db"
	"backend/server"
	"log"
)

func main() {
	database, err := db.Init()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Pass the database connection to StartServer
	server.StartServer(database)
}
