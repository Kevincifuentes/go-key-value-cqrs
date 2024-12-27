package config

type Config struct {
	Port        int    `env:"SERVER_PORT,required" envDefault:"8080"`
	OpenApiPath string `env:"OPENAPI_RELATIVE_PATH,required" envDefault:"./api/keyvalue/api.yml"`
	DebugMode   bool   `env:"DEBUG_MODE" envDefault:"false"`
}
