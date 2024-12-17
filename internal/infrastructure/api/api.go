package api

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yml ../../../api/keyvalue/api.yml

import (
	"github.com/getkin/kin-openapi/openapi3"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"net/http"
	"path/filepath"
)

type Server struct{}

func KeyValueServer() Server {
	return Server{}
}

func InitHandler(openApiRelativePath string) http.Handler {
	keyValueServer := KeyValueServer()
	serveMux := http.NewServeMux()

	handler := HandlerFromMux(keyValueServer, serveMux)
	openApiFilepath, _ := filepath.Abs(openApiRelativePath)
	swagger, _ := openapi3.NewLoader().LoadFromFile(openApiFilepath)
	validatorMiddleware := middleware.OapiRequestValidator(swagger)
	return validatorMiddleware(handler)
}

// GetKeyValueByKey (GET /keys/:key)
func (Server) GetKeyValueByKey(responseWriter http.ResponseWriter, _ *http.Request, _ string) {
	responseWriter.WriteHeader(http.StatusNotImplemented)
}

// PostKey (POST /keys)
func (Server) PostKey(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.WriteHeader(http.StatusNotImplemented)
}

// DeleteKeyValueByKey (DELETE /keys/{key})
func (Server) DeleteKeyValueByKey(responseWriter http.ResponseWriter, _ *http.Request, _ string) {
	responseWriter.WriteHeader(http.StatusNotImplemented)
}
