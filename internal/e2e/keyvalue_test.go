package e2e

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config config.yml ../../api/keyvalue/api.yml

import (
	"context"
	"go-key-value-cqrs/e2e/client"
	"go-key-value-cqrs/infrastructure/api"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/require"
)

var MaxKeyLength = 200
var FAKER = faker.New()

var testServer *httptest.Server
var keyValueClient *client.ClientWithResponses
var clientError error

func TestMain(testing *testing.M) {
	testServer = httptest.NewServer(api.InitHandler("../../api/keyvalue/api.yml"))

	hc := http.Client{}

	keyValueClient, clientError = client.NewClientWithResponses(testServer.URL, client.WithHTTPClient(&hc))
	if clientError != nil {
		log.Fatal(clientError)
	}

	testing.Run()
}

func TestGetKeyValueShouldNotBeFound(t *testing.T) {
	// given
	expectedKey := FAKER.Person().Name()

	// when
	response, err := keyValueClient.GetKeyValueByKeyWithResponse(context.TODO(), expectedKey)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, response.StatusCode())
	require.Equal(t, http.NoBody, response.HTTPResponse.Body)
}

func TestGetKeyValueBadRequestOnLongKey(t *testing.T) {
	// given
	expectedKey := FAKER.Lorem().Sentence(MaxKeyLength + 1)

	// when
	response, err := keyValueClient.GetKeyValueByKeyWithResponse(context.TODO(), expectedKey)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode())
}

func TestGetKeyValueFoundsValue(t *testing.T) {
	// given

	expectedKey := FAKER.Person().Name()
	expectedValue := FAKER.Person().Name()

	// when
	response, err := keyValueClient.GetKeyValueByKeyWithResponse(context.TODO(), expectedKey)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())
	require.Equal(t, map[string]string{expectedKey: expectedValue}, response.JSON200)
}
