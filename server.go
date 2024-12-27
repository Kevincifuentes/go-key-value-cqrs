package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"go-key-value-cqrs/infrastructure/api"
	"go-key-value-cqrs/infrastructure/api/config"
	"log"
	"net/http"
)

func main() {
	applicationConfig := retrieveConfiguration()
	handler := api.InitHandler(applicationConfig)

	server := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("0.0.0.0:%v", applicationConfig.Port),
	}
	log.Printf("Starting server on port %v\n", applicationConfig.Port)
	log.Println(server.ListenAndServe())
}

func retrieveConfiguration() config.Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.Config{}
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Unable to parse ennvironment variables: %e", err)
	}
	return cfg
}
