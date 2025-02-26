package main

import (
	"go-manage-mysql/cmd/config"
	"go-manage-mysql/internal/router"
	"log"
)

func main() {
	router := router.SetupRouter()
	if err := router.Run(config.Port); err != nil {
		log.Fatal("error starting server. Error: %w", err)
	}
}
