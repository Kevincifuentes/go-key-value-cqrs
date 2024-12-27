package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Port        int    `env:"SERVER_PORT,required" envDefault:"8080"`
	OpenApiPath string `env:"OPENAPI_RELATIVE_PATH,required" envDefault:"./api/keyvalue/api.yml"`
	DebugMode   bool   `env:"DEBUG_MODE" envDefault:"false"`
}

func RetrieveConfiguration() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Unable to parse ennvironment variables: %e", err)
	}
	return cfg
}
