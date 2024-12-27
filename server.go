package main

import (
	"fmt"
	"go-key-value-cqrs/infrastructure/api"
	"go-key-value-cqrs/infrastructure/api/config"
	"log"
	"net/http"
)

func main() {
	applicationConfig := config.RetrieveConfiguration()
	handler := api.InitHandler(applicationConfig)

	server := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("0.0.0.0:%v", applicationConfig.Port),
	}
	log.Printf("Starting server on port %v\n", applicationConfig.Port)
	log.Println(server.ListenAndServe())
}
