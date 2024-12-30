package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Host            string `env:"SERVER_HOST,required" envDefault:"localhost"`
	Port            int    `env:"SERVER_PORT,required" envDefault:"8080"`
	OpenApiPath     string `env:"OPENAPI_RELATIVE_PATH,required" envDefault:"./api.yml"`
	DebugServerPort int    `env:"DEBUG_SERVER_PORT" envDefault:"8081"`
	DebugServerHost string `env:"DEBUG_SERVER_HOST" envDefault:"localhost"`
}

func RetrieveConfiguration() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found. Using default environment variables...")
	}
	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Unable to parse ennvironment variables: %e", err)
	}
	return cfg
}

func (config Config) GetDebugServerAddress() string {
	return fmt.Sprintf("%s:%d", config.DebugServerHost, config.DebugServerPort)
}

func (config Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
