package api

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yml ../../../api/keyvalue/api.yml

import (
	"encoding/json"
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/gommon/log"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"go-key-value-cqrs/application/cqrs/commandbus"
	"go-key-value-cqrs/application/cqrs/querybus"
	"go-key-value-cqrs/application/keyvalue/addkeyvalue"
	"go-key-value-cqrs/application/keyvalue/deletekeyvalue"
	"go-key-value-cqrs/application/keyvalue/getvalue"
	"go-key-value-cqrs/domain"
	"go-key-value-cqrs/infrastructure/api/model"
	"go-key-value-cqrs/infrastructure/persistence"
	"golang.org/x/exp/maps"
	"net/http"
	"path/filepath"
)

type Server struct {
}

func InitHandler(openApiRelativePath string) http.Handler {
	keyValueServer := KeyValueServer()
	serveMux := http.NewServeMux()

	handler := HandlerFromMux(keyValueServer, serveMux)
	openApiFilepath, _ := filepath.Abs(openApiRelativePath)
	swagger, _ := openapi3.NewLoader().LoadFromFile(openApiFilepath)
	validatorMiddleware := middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{ErrorHandler: handleErrorMessage})
	return validatorMiddleware(handler)
}

func initializeQueryBus(keyValueReader *persistence.InMemoryKeyValueRepository) {
	err := querybus.Load(getvalue.QueryHandler{KeyValueReader: keyValueReader})
	if err != nil {
		log.Errorf("Error loading QueryBus %v", err)
	}
}

func initializeCommandBus(keyValueWriter *persistence.InMemoryKeyValueRepository) {
	var errs []error

	if err := commandbus.Load(addkeyvalue.CommandHandler{KeyValueWriter: keyValueWriter}); err != nil {
		errs = append(errs, err)
	}
	if err := commandbus.Load(deletekeyvalue.CommandHandler{KeyValueWriter: keyValueWriter}); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		log.Errorf("Error loading CommandBus: %v", errs)
	}
}

func KeyValueServer() Server {
	inMemoryKeyValueRepository := persistence.NewInMemoryKeyValueRepository()

	initializeQueryBus(inMemoryKeyValueRepository)
	initializeCommandBus(inMemoryKeyValueRepository)
	return Server{}
}

// GetKeyValueByKey (GET /keys/:key)
func (Server) GetKeyValueByKey(responseWriter http.ResponseWriter, _ *http.Request, key string) {
	keyValueView, err := querybus.Asks[domain.KeyValueView](getvalue.Query{Key: key})
	handleResponse(responseWriter, model.ToKeyValueResponse(keyValueView), http.StatusOK, err)
}

// PostKey (POST /keys)
func (Server) PostKey(responseWriter http.ResponseWriter, request *http.Request) {
	// Necessary because middleware not validating additionalProperties
	addKeyValue, key, err := validateAddKeyValueRequest(request)
	if err != nil {
		handleError(responseWriter, err)
		return
	}

	err = commandbus.Execute(addkeyvalue.Command{Key: key, Value: addKeyValue[key]})
	handleEmptyBodyResponse(responseWriter, http.StatusNoContent, err)
}

// DeleteKeyValueByKey (DELETE /keys/{key})
func (Server) DeleteKeyValueByKey(responseWriter http.ResponseWriter, _ *http.Request, key string) {
	err := commandbus.Execute(deletekeyvalue.Command{Key: key})
	handleEmptyBodyResponse(responseWriter, http.StatusNoContent, err)
}

func validateAddKeyValueRequest(request *http.Request) (AddKeyRequest, string, error) {
	var addKeyRequest AddKeyRequest
	if err := json.NewDecoder(request.Body).Decode(&addKeyRequest); err != nil {
		return nil, "", err
	}

	if len(addKeyRequest) != 1 {
		return nil, "", model.NewApiValidationError(addKeyRequest, "Only one key is allowed")
	}
	return addKeyRequest, maps.Keys(addKeyRequest)[0], nil
}

func handleResponse(writer http.ResponseWriter, response any, successfulStatusCode int, err error) {
	if err == nil {
		writeJSON(writer, successfulStatusCode, response)
		return
	}

	handleError(writer, err)
}

func handleError(writer http.ResponseWriter, err error) {
	var keyValueError *domain.KeyValueDomainError
	var apiValidationError model.ApiValidationError
	isKeyValueDomainError := errors.As(err, &keyValueError)
	isApiValidationError := errors.As(err, &apiValidationError)
	if isKeyValueDomainError || isApiValidationError {
		handleErrorMessage(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var keyNotFoundError *domain.KeyNotFoundError
	isKeyNotFoundError := errors.As(err, &keyNotFoundError)
	if isKeyNotFoundError {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	var keyAlreadyExists *domain.KeyExistsError
	isKeyAlreadyExists := errors.As(err, &keyAlreadyExists)
	if isKeyAlreadyExists {
		handleErrorMessage(writer, err.Error(), http.StatusConflict)
		return
	}

	handleErrorMessage(writer, err.Error(), http.StatusInternalServerError)
}

func handleErrorMessage(writer http.ResponseWriter, message string, statusCode int) {
	writeJSON(writer, statusCode, model.ErrorResponse{Message: message})
}

func handleEmptyBodyResponse(responseWriter http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		handleError(responseWriter, err)
		return
	}

	responseWriter.WriteHeader(statusCode)
}

func writeJSON(writer http.ResponseWriter, status int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	if err := json.NewEncoder(writer).Encode(data); err != nil {
		http.Error(writer, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
