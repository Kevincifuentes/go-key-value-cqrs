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

var maxKeyLength = 200
var fakerInstance = faker.New()

var testServer *httptest.Server
var keyValueClient *client.ClientWithResponses
var clientError error
var keyValueObjectMother KeyValueObjectMother

func TestMain(testing *testing.M) {
	testServer = httptest.NewServer(api.InitHandler("../../api/keyvalue/api.yml"))

	hc := http.Client{}

	keyValueClient, clientError = client.NewClientWithResponses(testServer.URL, client.WithHTTPClient(&hc))
	keyValueObjectMother = KeyValueObjectMother{fakerInstance: &fakerInstance, client: keyValueClient}
	if clientError != nil {
		log.Fatal(clientError)
	}

	testing.Run()
}

func TestGetKeyValueShouldNotBeFound(t *testing.T) {
	// given
	expectedKey := fakerInstance.Person().Name()

	// when
	response, err := keyValueClient.GetKeyValueByKeyWithResponse(context.TODO(), expectedKey)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, response.StatusCode())
	require.Equal(t, http.NoBody, response.HTTPResponse.Body)
}

func TestGetKeyValueBadRequestOnLongKey(t *testing.T) {
	// given
	expectedKey := fakerInstance.Lorem().Sentence(maxKeyLength + 1)

	// when
	response, err := keyValueClient.GetKeyValueByKeyWithResponse(context.TODO(), expectedKey)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode())
	require.Contains(t, *response.JSON400.Message, "maximum string length is 200")
}

func TestGetKeyValueFoundsValue(t *testing.T) {
	// given
	presentExpectedKey, presentExpectedValue := keyValueObjectMother.createRandom()

	// when
	response, err := keyValueClient.GetKeyValueByKeyWithResponse(context.TODO(), presentExpectedKey)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())
	expectedResponse := &client.KeyValueResponse{presentExpectedKey: presentExpectedValue}
	require.Equal(t, expectedResponse, response.JSON200)
}

func TestPostKeyValueSuccessful(t *testing.T) {
	// given
	expectedKey := fakerInstance.UUID().V4()
	expectedValue := fakerInstance.Person().Name()
	request := client.AddKeyRequest{expectedKey: expectedValue}

	// when
	response, err := keyValueClient.PostKeyWithResponse(context.TODO(), request)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode())
	require.Equal(t, http.NoBody, response.HTTPResponse.Body)
}

func TestPostKeyValueShouldReturnConflict(t *testing.T) {
	// given
	alreadyPresentExpectedKey, expectedValue := keyValueObjectMother.createRandom()
	request := client.AddKeyRequest{alreadyPresentExpectedKey: expectedValue}

	// when
	response, err := keyValueClient.PostKeyWithResponse(context.TODO(), request)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusConflict, response.StatusCode())
	require.Contains(t, string(response.Body), "already exists")
}

func TestPostKeyValueShouldReturnsBadRequestOnLongValues(t *testing.T) {
	// given
	expectedKey := fakerInstance.RandomStringWithLength(maxKeyLength + 1)
	expectedValue := fakerInstance.Person().Name()
	request := client.AddKeyRequest{expectedKey: expectedValue}

	// when
	response, err := keyValueClient.PostKeyWithResponse(context.TODO(), request)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode())
	require.Contains(t, string(response.Body), "expected 'key' to have a value between 1 and 200")
}

func TestPostKeyValueShouldReturnsBadRequestOnMoreThanOneKey(t *testing.T) {
	// given
	expectedKey := fakerInstance.UUID().V4()
	anotherKey := fakerInstance.UUID().V4()
	expectedValue := fakerInstance.Person().Name()
	request := client.AddKeyRequest{expectedKey: expectedValue, anotherKey: expectedValue}

	// when
	response, err := keyValueClient.PostKeyWithResponse(context.TODO(), request)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode())
	require.Contains(t, string(response.Body), "Only one key is allowed")
}
