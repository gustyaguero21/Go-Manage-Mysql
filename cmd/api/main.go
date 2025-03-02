package main

import (
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/database"
	"go-manage-mysql/internal/router"
	"log"
)

func main() {

	config.LoadEnv()

	conn, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	router := router.SetupRouter(conn)
	if err := router.Run(config.Port); err != nil {
		log.Fatal("error starting server. Error: %w", err)
	}
}
