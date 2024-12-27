package main

import (
	"go-key-value-cqrs/infrastructure/api"
	"go-key-value-cqrs/infrastructure/api/metrics"
	"log"
	"net/http"
	_ "net/http/pprof"
)

const OpenAPISpecRelativePath = "./api/keyvalue/api.yml"

func main() {
	handler := api.InitHandler(OpenAPISpecRelativePath)

	server := &http.Server{
		Handler: metrics.Metrics(handler),
		Addr:    "0.0.0.0:8080",
	}
	log.Println("Starting server on port 8080")
	log.Println(server.ListenAndServe())
}
