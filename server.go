package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-key-value-cqrs/infrastructure/api"
	"go-key-value-cqrs/infrastructure/api/metrics"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

const defaultOpenAPISpecRelativePath = "./api/keyvalue/api.yml"

func main() {
	loadEnvFile()

	openApiPath := retrieveOsVariableOrDefault("OPENAPI_RELATIVE_PATH", defaultOpenAPISpecRelativePath)
	handler := api.InitHandler(openApiPath)

	port := retrieveOsVariableOrDefault("PORT", "8080")
	server := &http.Server{
		Handler: metrics.Metrics(handler),
		Addr:    fmt.Sprintf("0.0.0.0:%v", port),
	}
	log.Printf("Starting server on port %v\n", port)
	log.Println(server.ListenAndServe())
}

func retrieveOsVariableOrDefault(environmentVariableName string, defaultValue string) string {
	relativePath, isPresent := os.LookupEnv(environmentVariableName)
	if !isPresent {
		return defaultValue
	}
	return relativePath
}

func loadEnvFile() {
	err := godotenv.Load(".env.test.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
