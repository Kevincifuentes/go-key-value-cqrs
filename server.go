package main

import (
	"go-key-value-cqrs/infrastructure/api"
	"log"
	"net/http"
)

const OpenAPISpecRelativePath = "./api/keyvalue/api.yml"

func main() {
	handler := api.InitHandler(OpenAPISpecRelativePath)
	server := http.Server{
		Handler: handler,
		Addr:    "0.0.0.0:8080",
	}
	log.Println("Starting server on port 8080")
	log.Fatal(server.ListenAndServe())
}
