package main

import (
	"gintraining/database"
	"gintraining/handlers"
	"gintraining/webapi"
	"log"
)

func main() {
	db := database.InitDB()
	h := handlers.InitHandler(db)

	srv := webapi.InitServer()
	err := srv.RunServer("8080", h.InitRoutes())
	if err != nil {
		log.Fatalf("server couldn't be started: %s", err)
	}
	log.Println("server has been started")
}
