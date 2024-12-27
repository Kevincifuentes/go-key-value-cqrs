package main

import (
	"go-key-value-cqrs/infrastructure/api"
	"go-key-value-cqrs/infrastructure/api/config"
	"log"
	"net/http"
)

func main() {
	applicationConfig := config.RetrieveConfiguration()
	handler := api.InitHandler(applicationConfig)

	serverAddress := applicationConfig.GetServerAddress()
	server := &http.Server{
		Handler: handler,
		Addr:    serverAddress,
	}
	log.Printf("Starting KeyValueServer on %v\n", serverAddress)
	log.Println(server.ListenAndServe())
}
