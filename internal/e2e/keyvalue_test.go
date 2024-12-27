package e2e

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config config.yml ../../api/keyvalue/api.yml

import (
	"context"
	"github.com/joho/godotenv"
	"go-key-value-cqrs/e2e/client"
	"go-key-value-cqrs/infrastructure/api"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
	loadEnvFile()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	var hostUrl string
	if loadTestServer, _ := strconv.ParseBool(os.Getenv("LOAD_TEST_SERVER")); !loadTestServer {
		testServer = httptest.NewServer(api.InitHandler(os.Getenv("OPENAPI_RELATIVE_PATH")))
		hostUrl = testServer.URL
	} else {
		hostUrl = os.Getenv("LOCAL_HOST_URL")
	}

	log.Printf("Targeting server on %s...\n", hostUrl)

	hc := http.Client{}
	keyValueClient, clientError = client.NewClientWithResponses(hostUrl, client.WithHTTPClient(&hc))
	keyValueObjectMother = KeyValueObjectMother{fakerInstance: &fakerInstance, client: keyValueClient}
	if clientError != nil {
		log.Fatal(clientError)
	}

	testing.Run()
}

func loadEnvFile() {
	err := godotenv.Load(".env.test.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestGetKeyValueShouldNotBeFound(t *testing.T) {
	// given
	unknownKey := fakerInstance.Person().Name()

	// when
	response, err := keyValueClient.GetKeyValueByKeyWithResponse(context.TODO(), unknownKey)

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

func TestDeleteKeyValueShouldBeSuccessful(t *testing.T) {
	// given
	alreadyPresentExpectedKey, _ := keyValueObjectMother.createRandom()

	// when
	response, err := keyValueClient.DeleteKeyValueByKeyWithResponse(context.TODO(), alreadyPresentExpectedKey)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode())
	require.Equal(t, http.NoBody, response.HTTPResponse.Body)
}

func TestDeleteKeyValueShouldReturnBadRequestOnLongKey(t *testing.T) {
	// given
	tooLongKey := fakerInstance.RandomStringWithLength(maxKeyLength + 1)

	// when
	response, err := keyValueClient.DeleteKeyValueByKeyWithResponse(context.TODO(), tooLongKey)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode())
	require.Contains(t, string(response.Body), "maximum string length is 200")
}

func TestDeleteKeyValueShouldReturnNotFoundOnMissingKey(t *testing.T) {
	// given
	unknownKey := fakerInstance.UUID().V4()

	// when
	response, err := keyValueClient.DeleteKeyValueByKeyWithResponse(context.TODO(), unknownKey)

	// then
	require.Nil(t, err)
	require.Equal(t, http.StatusNotFound, response.StatusCode())
	require.Equal(t, http.NoBody, response.HTTPResponse.Body)
}

func BenchmarkGetKeyValueWithRandomKeys(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// given
			unknownKey := fakerInstance.UUID().V4()

			// when
			_, _ = keyValueClient.GetKeyValueByKeyWithResponse(context.TODO(), unknownKey)
		}
	})
}

func BenchmarkPostKeyValueWithRandomKeys(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// when
			keyValueObjectMother.createRandom()
		}
	})
}
