package api

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yml ../../../api/keyvalue/api.yml

import (
	"encoding/json"
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/gommon/log"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"go-key-value-cqrs/application/queries/cqrs/querybus"
	"go-key-value-cqrs/application/queries/keyvalue/getvalue"
	"go-key-value-cqrs/domain"
	"go-key-value-cqrs/infrastructure/api/model"
	"go-key-value-cqrs/infrastructure/persistence"
	"net/http"
	"path/filepath"
)

type Server struct {
	queryBus querybus.Query
}

func initializeQueryBus() {
	keyValueReader := persistence.NewInMemoryKeyValueRepository()

	err := querybus.Load(getvalue.QueryHandler{KeyValueReader: keyValueReader})
	if err != nil {
		log.Errorf("Error loading query bus %s", err)
	}
}

func KeyValueServer() Server {
	initializeQueryBus()
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
func (Server) GetKeyValueByKey(responseWriter http.ResponseWriter, _ *http.Request, key string) {
	keyValueView, err := querybus.Asks[domain.KeyValueView](getvalue.Query{Key: key})
	handleResponse(responseWriter, model.ToKeyValueResponse(keyValueView), http.StatusOK, err)
}

// PostKey (POST /keys)
func (Server) PostKey(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.WriteHeader(http.StatusNotImplemented)
}

// DeleteKeyValueByKey (DELETE /keys/{key})
func (Server) DeleteKeyValueByKey(responseWriter http.ResponseWriter, _ *http.Request, _ string) {
	responseWriter.WriteHeader(http.StatusNotImplemented)
}

func handleResponse(writer http.ResponseWriter, response any, successfulStatusCode int, err error) {
	if err == nil {
		writeJSON(writer, successfulStatusCode, response)
		return
	}

	var keyValueError *domain.KeyValueDomainError
	isKeyValueDomainError := errors.As(err, &keyValueError)
	if isKeyValueDomainError {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var keyNotFoundError *domain.KeyNotFoundError
	isKeyNotFoundError := errors.As(err, &keyNotFoundError)
	if isKeyNotFoundError {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	http.Error(writer, err.Error(), http.StatusInternalServerError)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
