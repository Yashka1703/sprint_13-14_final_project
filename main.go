package main

import (
	"fmt"
	"log"

	"finalProject/pkg/api"
	"finalProject/pkg/db"
	"finalProject/pkg/server"
)

func main() {

	dbFile := "scheduler.db"

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("db init error %v", err)
	}

	defer func() {
		if err := db.GetDB().Close(); err != nil {
			log.Printf("error when closing DB: %v", err)
		}
	}()
	api.Init()

	log.Println("Starting server")
	if err := server.StartServ(); err != nil {
		fmt.Printf("Error starting server: %v", err)
	}

}
